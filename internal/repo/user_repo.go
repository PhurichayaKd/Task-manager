package repo

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"task-manager/internal/domain"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

type UserRepo interface {
	// สมัครแบบ local (มีรหัสผ่าน)
	Create(ctx context.Context, u *domain.User) (*domain.User, error)

	// สมัคร/ล็อกอินผ่าน OAuth
	CreateFromOAuth(ctx context.Context, u *domain.User) (int64, error)

	// ดึงด้วยอีเมล
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	// ดึงด้วย username
	GetByUsername(ctx context.Context, username string) (*domain.User, error)

	// ดึงด้วย ID
	GetByID(ctx context.Context, id int64) (*domain.User, error)

	// อัปเดต user
	Update(ctx context.Context, u *domain.User) (*domain.User, error)

	// อัปเดตชื่อ
	UpdateName(ctx context.Context, id int64, name string) error

	// อัปเดต username
	UpdateUsername(ctx context.Context, id int64, username string) error

	// อัปเดตรหัสผ่าน
	UpdatePassword(ctx context.Context, id int64, hashedPassword string) error

	// เช็กว่ามีอีเมลนี้หรือยัง
	EmailExists(ctx context.Context, email string) (bool, error)

	// เช็กว่ามี username นี้หรือยัง
	UsernameExists(ctx context.Context, username string) (bool, error)
}

type userRepo struct{ db *sql.DB }

func NewUserRepo(db *sql.DB) UserRepo { return &userRepo{db: db} }

func MustOpen(dsn string) *sql.DB {
	var driver string
	if strings.HasPrefix(dsn, "sqlite3://") {
		driver = "sqlite"
		dsn = strings.TrimPrefix(dsn, "sqlite3://")
	} else if strings.HasPrefix(dsn, "file:") {
		driver = "sqlite"
	} else {
		driver = "postgres"
	}
	
	db, err := sql.Open(driver, dsn)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}
	return db
}

func (r *userRepo) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx,
		`INSERT INTO users (email, username, password_hash, role)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, email, username, password_hash, role, created_at`,
		u.Email, u.Username, u.PasswordHash, u.Role,
	)

	var out domain.User
	if err := row.Scan(&out.ID, &out.Email, &out.Username, &out.PasswordHash, &out.Role, &out.CreatedAt); err != nil {
		return nil, err
	}
	return &out, nil
}

// ใช้สำหรับสมัครด้วย Google OAuth (ไม่มี password)
func (r *userRepo) CreateFromOAuth(ctx context.Context, u *domain.User) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// หมายเหตุ: ฟิลด์ให้ตรงกับคอลัมน์จริงใน DB ของคุณ
	// แนะนำให้มี role ค่า default เป็น 'user'
	var id int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO users (email, name, provider, provider_id, avatar_url, role)
		VALUES ($1,$2,$3,$4,$5, COALESCE($6,'user'))
		RETURNING id
	`, u.Email, u.Name, u.Provider, u.ProviderID, u.AvatarURL, u.Role).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, `
		SELECT id, email, username, password_hash, role,
		       name, provider, provider_id, avatar_url,
		       created_at
		FROM users
		WHERE email = $1
	`, email)

	var u domain.User
	if err := row.Scan(
		&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.Role,
		&u.Name, &u.Provider, &u.ProviderID, &u.AvatarURL,
		&u.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, `
		SELECT id, email, username, password_hash, role,
		       name, provider, provider_id, avatar_url,
		       created_at
		FROM users
		WHERE username = $1
	`, username)

	var u domain.User
	if err := row.Scan(
		&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.Role,
		&u.Name, &u.Provider, &u.ProviderID, &u.AvatarURL,
		&u.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) EmailExists(ctx context.Context, email string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var ok bool
	if err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, email).Scan(&ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (r *userRepo) UsernameExists(ctx context.Context, username string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var ok bool
	if err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`, username).Scan(&ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, email, name, created_at FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var u domain.User
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *userRepo) UpdateName(ctx context.Context, id int64, name string) error {
	query := `UPDATE users SET name = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, name, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *userRepo) UpdateUsername(ctx context.Context, id int64, username string) error {
	query := `UPDATE users SET username = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, username, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *userRepo) UpdatePassword(ctx context.Context, id int64, hashedPassword string) error {
	query := `UPDATE users SET password_hash = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, hashedPassword, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *userRepo) Update(ctx context.Context, u *domain.User) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		UPDATE users 
		SET username = $1, password_hash = $2, name = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query, u.Username, u.PasswordHash, u.Name, u.ID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrNotFound
	}

	// Return updated user
	return r.GetByID(ctx, u.ID)
}

var ErrNotFound = errors.New("not found")
