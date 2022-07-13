accrualwin:
	./cmd/accrual/accrual_windows_amd64

accruals:
	curl.exe -X POST http://localhost:8080/api/goods -H "Content-Type: application/json" -d '{ "match": "LG", "reward": 5, "reward_type": "%" }'
	curl.exe -X POST http://localhost:8080/api/goods -H "Content-Type: application/json" -d '{ "match": "Nokia", "reward": 2, "reward_type": "%" }'
	curl.exe -X POST http://localhost:8080/api/goods -H "Content-Type: application/json" -d '{ "match": "Festool", "reward": 10, "reward_type": "%" }'
	curl.exe -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d '{ "order": "123455", "goods": [ { "description": "LG Monitor", "price": 50000.0 } ] }'
	curl.exe -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d '{ "order": "2675865373825070", "goods": [ { "description": "Festool Saw", "price": 20000.0 } ] }'
	curl.exe -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d '{ "order": "6188644838072821", "goods": [ { "description": "Nokia 3310", "price": 30000.0 } ] }'
	curl.exe -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d '{ "order": "2254083172131232", "goods": [ { "description": "Festool Saw", "price": 50000.0 } ] }'
	curl.exe -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d '{ "order": "4225747548588380", "goods": [ { "description": "LG Monitor", "price": 60000.0 } ] }'
	curl.exe -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d '{ "order": "0740064321441447", "goods": [ { "description": "Festool Saw", "price": 80000.0 } ] }'
	curl.exe -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d '{ "order": "3663733635151572", "goods": [ { "description": "LG Monitor", "price": 20000.0 } ] }'
	curl.exe -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d '{ "order": "1701803731677230", "goods": [ { "description": "Nokia 3310", "price": 40000.0 } ] }'
	curl.exe -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d '{ "order": "4246720124876623", "goods": [ { "description": "LG Monitor", "price": 60000.0 } ] }'
	curl.exe -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d '{ "order": "1258213818211688", "goods": [ { "description": "Festool Saw", "price": 50000.0 } ] }'
	curl.exe -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d '{ "order": "2044515743876311", "goods": [ { "description": "LG Monitor", "price": 60000.0 } ] }'
	curl.exe -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d '{ "order": "6103701024561769", "goods": [ { "description": "Nokia 3310", "price": 50000.0 } ] }'
	curl.exe -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d '{ "order": "2141885384266300", "goods": [ { "description": "Festool Saw", "price": 40000.0 } ] }'
	curl.exe -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d '{ "order": "4216744344110106", "goods": [ { "description": "LG Monitor", "price": 50000.0 } ] }'
