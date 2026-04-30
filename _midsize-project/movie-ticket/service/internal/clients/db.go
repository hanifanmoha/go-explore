package clients

import (
	"context"
	"fmt"

	"github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseClient struct {
	Pool *pgxpool.Pool
}

func NewDatabaseClient() (*DatabaseClient, error) {
	dbUrl := config.GetDbUrl()
	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	// Initialize database schema if not exists
	// In production, you should use proper migration tools like golang-migrate or goose
	if err := UnproperlyInitDatabaseSchema(pool); err != nil {
		return nil, err
	}

	return &DatabaseClient{Pool: pool}, nil
}

func (c *DatabaseClient) Close() {
	c.Pool.Close()
}

func UnproperlyInitDatabaseSchema(pool *pgxpool.Pool) error {

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS movies (
		id SERIAL PRIMARY KEY,
		title TEXT UNIQUE NOT NULL,
		description TEXT,
		image_url TEXT,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	INSERT INTO movies (title, description, image_url) VALUES
	('The Bourne Identity', 'A man is picked up by a fishing boat, bullet-riddled and suffering from amnesia, before racing to elude assassins and attempting to regain his memory.', 'https://m.media-amazon.com/images/M/MV5BYTk1ZTcyMWMtMWUxYS00MmEzLTlmODYtOTk1MGRjOTg1ZjlmXkEyXkFqcGc@._V1_FMjpg_UX1000_.jpg'),
	('The Bourne Supremacy', 'When Jason Bourne is framed for a CIA operation gone awry, he is forced to resume his former life as a trained assassin to survive.', 'https://m.media-amazon.com/images/M/MV5BZTU4ZDgyYjgtODA0Mi00MmE3LTgzYWQtZjc1YTFiMTczZTQ3XkEyXkFqcGc@._V1_FMjpg_UX1014_.jpg'),
	('The Bourne Ultimatum', 'Jason Bourne dodges a ruthless C.I.A. official and his Agents from a new assassination program while searching for the origins of his life as a trained killer.', 'https://m.media-amazon.com/images/M/MV5BYzE3ZGM4MzctZjU5MC00NWE2LTg5ZjYtMDFiM2ZlMWQ1MjkwXkEyXkFqcGc@._V1_FMjpg_UY3000_.jpg')
	ON CONFLICT (title) DO NOTHING;

	CREATE TABLE IF NOT EXISTS seats (
		id SERIAL PRIMARY KEY,
		movie_id INTEGER NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
		seat_number VARCHAR(10) NOT NULL,
		status TEXT NOT NULL DEFAULT 'available',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		UNIQUE(movie_id, seat_number)
	);

	CREATE TABLE IF NOT EXISTS bookings (
		id SERIAL PRIMARY KEY,
		movie_id INTEGER NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
		seat_number VARCHAR(10) NOT NULL,
		user_id INTEGER NOT NULL,
		status TEXT NOT NULL DEFAULT 'booked',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		UNIQUE(movie_id, seat_number)
	);
	`

	_, err := pool.Exec(context.Background(), createTableQuery)
	if err != nil {
		return err
	}

	rows := []string{"A", "B", "C", "D"}
	colNumber := 8
	insertSeatsQuery := ``
	for _, row := range rows {
		for col := 1; col <= colNumber; col++ {
			seatNumber := fmt.Sprintf("%s%d", row, col)
			insertSeatsQuery += `INSERT INTO seats (movie_id, seat_number) VALUES (1, '` + seatNumber + `') ON CONFLICT (movie_id, seat_number) DO NOTHING;`
			insertSeatsQuery += `INSERT INTO seats (movie_id, seat_number) VALUES (2, '` + seatNumber + `') ON CONFLICT (movie_id, seat_number) DO NOTHING;`
			insertSeatsQuery += `INSERT INTO seats (movie_id, seat_number) VALUES (3, '` + seatNumber + `') ON CONFLICT (movie_id, seat_number) DO NOTHING;`
		}
	}

	_, err = pool.Exec(context.Background(), insertSeatsQuery)
	if err != nil {
		return err
	}

	return nil
}
