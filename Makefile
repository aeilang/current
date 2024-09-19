migration_up: migration_clean
	@migrate -database="postgres://lang:password@localhost:5432/test_db?sslmode=disable" \
	-path="migration" up


migration_clean:
	@migrate -database="postgres://lang:password@localhost:5432/test_db?sslmode=disable" \
	-path="migration" drop -f

curl2:
	curl -X PUT -d '{"age": 19}' http://localhost:8888/student/1 & \
	curl -X PUT -d '{"name": "李四"}' http://localhost:8888/student/1 & \
	wait