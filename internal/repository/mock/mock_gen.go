package mock

//go:generate mockgen -source=../../domain/repository/user_repository.go -destination=user_repository_mock.go -package=mock
//go:generate mockgen -source=../../domain/repository/pvz_repository.go -destination=pvz_repository_mock.go -package=mock
//go:generate mockgen -source=../../domain/repository/reception_repository.go -destination=reception_repository_mock.go -package=mock
//go:generate mockgen -source=../../domain/repository/product_repository.go -destination=product_repository_mock.go -package=mock
