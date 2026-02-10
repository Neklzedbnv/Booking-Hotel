package db

import (
	"database/sql"
	"log"
)

// RunMigrations создаёт все таблицы если их ещё нет.
// Вызывайте один раз при старте приложения.
func RunMigrations(db *sql.DB) {
	queries := []string{
		// ── users ────────────────────────────────────────────────────
		`CREATE TABLE IF NOT EXISTS users (
			id            SERIAL PRIMARY KEY,
			fullname      VARCHAR(255) NOT NULL,
			email         VARCHAR(255) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			role          VARCHAR(50)  NOT NULL DEFAULT 'user',
			created_at    TIMESTAMP    NOT NULL DEFAULT NOW()
		);`,

		// ── room_types ──────────────────────────────────────────────
		`CREATE TABLE IF NOT EXISTS room_types (
			id         BIGSERIAL PRIMARY KEY,
			name       VARCHAR(100) NOT NULL UNIQUE,
			capacity   INT          NOT NULL,
			base_price NUMERIC(10,2) NOT NULL,
			created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
		);`,

		// ── rooms ───────────────────────────────────────────────────
		`CREATE TABLE IF NOT EXISTS rooms (
			id         BIGSERIAL PRIMARY KEY,
			code       VARCHAR(50)   NOT NULL UNIQUE,
			type_id    BIGINT        NOT NULL REFERENCES room_types(id) ON DELETE CASCADE,
			capacity   INT           NOT NULL,
			price      NUMERIC(10,2) NOT NULL,
			status     VARCHAR(50)   NOT NULL DEFAULT 'available',
			created_at TIMESTAMP     NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP     NOT NULL DEFAULT NOW()
		);`,

		// ── meal_plans ──────────────────────────────────────────────
		`CREATE TABLE IF NOT EXISTS meal_plans (
			id            BIGSERIAL PRIMARY KEY,
			name          VARCHAR(100)  NOT NULL UNIQUE,
			price_per_day NUMERIC(10,2) NOT NULL,
			created_at    TIMESTAMP     NOT NULL DEFAULT NOW()
		);`,

		// ── packages ────────────────────────────────────────────────
		`CREATE TABLE IF NOT EXISTS packages (
			id             SERIAL PRIMARY KEY,
			name           VARCHAR(255)  NOT NULL,
			description    TEXT          NOT NULL DEFAULT '',
			price_modifier NUMERIC(10,2) NOT NULL DEFAULT 0,
			is_active      BOOLEAN       NOT NULL DEFAULT TRUE,
			created_at     TIMESTAMP     NOT NULL DEFAULT NOW(),
			updated_at     TIMESTAMP     NOT NULL DEFAULT NOW()
		);`,

		// ── room_packages (many-to-many) ────────────────────────────
		`CREATE TABLE IF NOT EXISTS room_packages (
			id         BIGSERIAL PRIMARY KEY,
			room_id    BIGINT    NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
			package_id INT       NOT NULL REFERENCES packages(id) ON DELETE CASCADE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE(room_id, package_id)
		);`,

		// ── bookings ────────────────────────────────────────────────
		`CREATE TABLE IF NOT EXISTS bookings (
			id          SERIAL PRIMARY KEY,
			user_id     INT           NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			room_id     INT           NOT NULL,
			mealplan_id INT           REFERENCES meal_plans(id) ON DELETE SET NULL,
			package_id  INT           REFERENCES packages(id) ON DELETE SET NULL,
			start_date  DATE          NOT NULL,
			end_date    DATE          NOT NULL,
			stay_days   INT           NOT NULL DEFAULT 0,
			total_price NUMERIC(10,2) NOT NULL DEFAULT 0,
			created_at  TIMESTAMP     NOT NULL DEFAULT NOW()
		);`,

		// ── payments ────────────────────────────────────────────────
		`CREATE TABLE IF NOT EXISTS payments (
			id         SERIAL PRIMARY KEY,
			booking_id INT           NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
			method     VARCHAR(50)   NOT NULL,
			status     VARCHAR(50)   NOT NULL DEFAULT 'pending',
			amount     NUMERIC(10,2) NOT NULL,
			created_at TIMESTAMP     NOT NULL DEFAULT NOW()
		);`,

		// ── reviews ─────────────────────────────────────────────────
		`CREATE TABLE IF NOT EXISTS reviews (
			id         SERIAL PRIMARY KEY,
			booking_id INT       NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
			rating     INT       NOT NULL CHECK (rating >= 1 AND rating <= 5),
			comment    TEXT      NOT NULL DEFAULT '',
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);`,

		// ── services ────────────────────────────────────────────────
		`CREATE TABLE IF NOT EXISTS services (
			id         SERIAL PRIMARY KEY,
			name       VARCHAR(255)  NOT NULL,
			price      NUMERIC(10,2) NOT NULL,
			created_at TIMESTAMP     NOT NULL DEFAULT NOW()
		);`,

		// ── indexes ─────────────────────────────────────────────────
		`CREATE INDEX IF NOT EXISTS idx_users_email        ON users(email);`,
		`CREATE INDEX IF NOT EXISTS idx_rooms_type_id      ON rooms(type_id);`,
		`CREATE INDEX IF NOT EXISTS idx_rooms_status       ON rooms(status);`,
		`CREATE INDEX IF NOT EXISTS idx_bookings_user_id   ON bookings(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_bookings_room_id   ON bookings(room_id);`,
		`CREATE INDEX IF NOT EXISTS idx_bookings_dates     ON bookings(start_date, end_date);`,
		`CREATE INDEX IF NOT EXISTS idx_payments_booking   ON payments(booking_id);`,
		`CREATE INDEX IF NOT EXISTS idx_reviews_booking    ON reviews(booking_id);`,
		`CREATE INDEX IF NOT EXISTS idx_room_packages_room ON room_packages(room_id);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Fatalf("migration error: %v\nquery: %s", err, q)
		}
	}

	log.Println("migrations completed successfully")
}
