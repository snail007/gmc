module github.com/snail007/gmc

go 1.12

require (
	github.com/dsnet/compress v0.0.1
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.1.2
	github.com/klauspost/pgzip v1.2.5
	github.com/pkg/errors v0.8.1
	github.com/snail007/gmct v0.0.23
	github.com/snail007/go-sqlcipher v0.0.0-20210114093415-fb27975e042f
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.5.0
	github.com/ulikunitz/xz v0.5.8
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b
	golang.org/x/text v0.3.3
)

replace github.com/snail007/gmct => ../gmct
