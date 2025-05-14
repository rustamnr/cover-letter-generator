NGROK_AUTH_TOKEN=2x0bvTOdJS1ZyWGfTy3PuJYc2p2_6KRasKnNBARGFruZX41Np

# ===== Start application =====
.PHONY: run

run:
	make redis
	make ngrok 
	go run ./cmd/main.go

# ===== Start ngrok and update API_BASE_URL in .env =====

.PHONY: ngrok ngrok-url update-env stop-ngrok

ngrok: ngrok-auth stop-ngrok ngrok-url update-env

ngrok-auth:
	@echo "Authenticating ngrok..."
	@ngrok config add-authtoken $(NGROK_AUTH_TOKEN)

ngrok-url:
	@NGROK_URL=$$(cat .ngrok_url); \
	if grep -q '^API_BASE_URL=' .env; then \
		sed -i "s|^API_BASE_URL=.*|API_BASE_URL=$$NGROK_URL|" .env; \
	else \
		echo "API_BASE_URL=$$NGROK_URL" >> .env; \
	fi; \
	echo "API_BASE_URL updated: $$NGROK_URL"

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
