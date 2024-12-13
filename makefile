FRONTEND_DIR = frontend
WAILS_WEB_PORT = 34115
E2E_DB_PATH = ./frontend/tests/db

dev:
	wails dev -tags webkit2_41

package:
	wails build -tags webkit2_41

goose_up:
	goose -dir=internal/database/migrations sqlite3 ~/.config/progressor/progressor.db up

goose_down:
	goose -dir=internal/database/migrations sqlite3 ~/.config/progressor/progressor.db down

start:
	@echo "Starting Wails application..."
	rm -f ${E2E_DB_PATH}/progressor.db
	DATABASE_PATH=$(E2E_DB_PATH) wails dev -tags webkit2_41 & \
	PID=$$!; \
	echo $$PID > wails.pid; \
	while ! nc -z localhost $(WAILS_WEB_PORT); do sleep 1; done; \
	echo "Wails app is running on port $(WAILS_WEB_PORT)"

# Target to run Playwright tests
test:
	@echo "Running Playwright tests..."
	cd $(FRONTEND_DIR) && npx playwright test

# Target to clean up
clean:
	@echo "Cleaning up..."
	rm -f ${E2E_DB_PATH}/progressor.db
	@kill $$(cat wails.pid) || true; \
	rm -f wails.pid

# Combined target to build, start, and test
test_e2e: start test clean
