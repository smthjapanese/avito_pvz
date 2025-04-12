package mock

//go:generate mockgen -source=../../usecase/user_usecase.go -destination=mock_user_usecase.go -package=mock
//go:generate mockgen -source=../../usecase/pvz_usecase.go -destination=mock_pvz_usecase.go -package=mock
//go:generate mockgen -source=../../usecase/reception_usecase.go -destination=mock_reception_usecase.go -package=mock
//go:generate mockgen -source=../../usecase/product_usecase.go -destination=mock_product_usecase.go -package=mock
