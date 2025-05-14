# ===== Start application =====
.PHONY: run

run:
	make redis
	make ngrok 
	go run ./cmd/main.go

# ===== Start ngrok and update API_BASE_URL in .env =====

.PHONY: ngrok ngrok-url update-env stop-ngrok

ngrok: stop-ngrok ngrok-url update-env

ngrok-url:
	@echo "Start ngrok..."
	@ngrok http 8080 > /dev/null & \
	sleep 3 && \
	curl -s http://localhost:4040/api/tunnels | \
	grep -o '"public_url":"https:[^"]*' | \
	sed 's/"public_url":"//' > .ngrok_url

update-env:
	@NGROK_URL=$$(cat .ngrok_url); \
	if grep -q '^API_BASE_URL=' .env; then \
		sed -i '' "s|^API_BASE_URL=.*|API_BASE_URL=$$NGROK_URL|" .env; \
	else \
		echo "API_BASE_URL=$$NGROK_URL" >> .env; \
	fi; \
	echo "API_BASE_URL updated: $$NGROK_URL"

stop-ngrok:
	@pkill ngrok || true

# ===== Start Redis container =====

.PHONY: redis

redis:
	@if [ -z "$$(docker ps -q -f name=coverletter-redis)" ]; then \
		if [ -z "$$(docker ps -aq -f name=coverletter-redis)" ]; then \
			echo "Creating Redis container..."; \
			docker run -d --name coverletter-redis -p 6379:6379 redis:7-alpine; \
		else \
			echo "Start Redis docker-container ..."; \
			docker start coverletter-redis; \
		fi \
	else \
		echo "Redis is already running."; \
	fi
