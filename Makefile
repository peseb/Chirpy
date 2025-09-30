migrate: 
	cd sql/schema && goose postgres "postgres://petbjo:@localhost:5432/chirpy" up