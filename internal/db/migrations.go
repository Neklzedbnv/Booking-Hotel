package db

import (
	"database/sql"
	"log"
)

// RunMigrations creates all tables if they don't exist yet.
// Call once at application startup.
func RunMigrations(db *sql.DB) {
	queries := []string{
		// ── users ────────────────────────────────────────────────────
		`CREATE TABLE IF NOT EXISTS users (
			id            SERIAL PRIMARY KEY,
			fullname      VARCHAR(255) NOT NULL,
			email         VARCHAR(255) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			role          VARCHAR(50)  NOT NULL DEFAULT 'user',
			is_blocked    BOOLEAN      NOT NULL DEFAULT FALSE,
			created_at    TIMESTAMP    NOT NULL DEFAULT NOW()
		);`,

		// add is_blocked column if missing
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS is_blocked BOOLEAN NOT NULL DEFAULT FALSE;`,

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
			id             BIGSERIAL PRIMARY KEY,
			code           VARCHAR(50)   NOT NULL UNIQUE,
			room_type_id   BIGINT        NOT NULL REFERENCES room_types(id) ON DELETE CASCADE,
			capacity       INT           NOT NULL,
			price          NUMERIC(10,2) NOT NULL,
			status         VARCHAR(50)   NOT NULL DEFAULT 'available',
			created_at     TIMESTAMP     NOT NULL DEFAULT NOW(),
			updated_at     TIMESTAMP     NOT NULL DEFAULT NOW()
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
			status      VARCHAR(50)   NOT NULL DEFAULT 'pending',
			created_at  TIMESTAMP     NOT NULL DEFAULT NOW()
		);`,

		// add status column if missing
		`ALTER TABLE bookings ADD COLUMN IF NOT EXISTS status VARCHAR(50) NOT NULL DEFAULT 'pending';`,

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
		`CREATE INDEX IF NOT EXISTS idx_rooms_type_id      ON rooms(room_type_id);`,
		`CREATE INDEX IF NOT EXISTS idx_rooms_status       ON rooms(status);`,
		`CREATE INDEX IF NOT EXISTS idx_bookings_user_id   ON bookings(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_bookings_room_id   ON bookings(room_id);`,
		`CREATE INDEX IF NOT EXISTS idx_bookings_dates     ON bookings(start_date, end_date);`,
		`CREATE INDEX IF NOT EXISTS idx_payments_booking   ON payments(booking_id);`,
		`CREATE INDEX IF NOT EXISTS idx_reviews_booking    ON reviews(booking_id);`,
		`CREATE INDEX IF NOT EXISTS idx_room_packages_room ON room_packages(room_id);`,

		// ── Translate Russian room type names to English ────────────
		`UPDATE room_types SET name = 'Standard' WHERE name = 'Стандарт';`,
		`UPDATE room_types SET name = 'Deluxe' WHERE name = 'Делюкс';`,
		`UPDATE room_types SET name = 'Suite' WHERE name = 'Люкс';`,
		`UPDATE room_types SET name = 'Premium' WHERE name = 'Премиум';`,
		`UPDATE room_types SET name = 'Business' WHERE name = 'Бизнес';`,
		`UPDATE room_types SET name = 'Family' WHERE name = 'Семейный';`,
		`UPDATE room_types SET name = 'Economy' WHERE name = 'Эконом';`,

		// ── Translate Russian service names to English ──────────────
		`UPDATE services SET name = 'Spa Treatments' WHERE name = 'СПА процедуры' OR name = 'Спа процедуры' OR name = 'СПА';`,
		`UPDATE services SET name = 'Gym' WHERE name = 'Тренажерный зал' OR name = 'Спортзал' OR name = 'Фитнес';`,
		`UPDATE services SET name = 'Bike Rental' WHERE name = 'Аренда велосипеда' OR name = 'Прокат велосипедов';`,
		`UPDATE services SET name = 'Tours' WHERE name = 'Экскурсии' OR name = 'Туры';`,
		`UPDATE services SET name = 'Restaurant' WHERE name = 'Ресторан';`,
		`UPDATE services SET name = 'Breakfast' WHERE name = 'Завтрак';`,
		`UPDATE services SET name = 'Laundry' WHERE name = 'Прачечная' OR name = 'Стирка';`,
		`UPDATE services SET name = 'Parking' WHERE name = 'Парковка';`,
		`UPDATE services SET name = 'Transfer' WHERE name = 'Трансфер';`,
		`UPDATE services SET name = 'Mini Bar' WHERE name = 'Мини-бар' OR name = 'Минибар';`,

		// ── Translate Russian meal plan names to English ────────────
		`UPDATE meal_plans SET name = 'Breakfast Only' WHERE name = 'Только завтрак' OR name = 'Завтрак';`,
		`UPDATE meal_plans SET name = 'Half Board' WHERE name = 'Полупансион';`,
		`UPDATE meal_plans SET name = 'Full Board' WHERE name = 'Полный пансион';`,
		`UPDATE meal_plans SET name = 'All Inclusive' WHERE name = 'Все включено' OR name = 'Всё включено';`,
		`UPDATE meal_plans SET name = 'No Meals' WHERE name = 'Без питания';`,

		// ── Translate Russian package names to English ──────────────
		`UPDATE packages SET name = 'Romantic' WHERE name = 'Романтический';`,
		`UPDATE packages SET name = 'Business' WHERE name = 'Бизнес';`,
		`UPDATE packages SET name = 'Family' WHERE name = 'Семейный';`,
		`UPDATE packages SET name = 'Relax' WHERE name = 'Релакс' OR name = 'Отдых';`,
		`UPDATE packages SET name = 'Weekend' WHERE name = 'Выходной' OR name = 'Выходные';`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Printf("migration warning: %v", err)
		}
	}

	log.Println("migrations completed")
}
