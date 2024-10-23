package main

import (
	"github.com/joho/godotenv"
	"github.com/minhnghia2k3/exchanger/internal/database"
	"github.com/minhnghia2k3/exchanger/internal/env"
	"github.com/minhnghia2k3/exchanger/internal/mail"
	"github.com/minhnghia2k3/exchanger/internal/store"
	"log"
	"log/slog"
	"os"
)

const jsonData = `
{
    "result": "success",
    "documentation": "https://www.exchangerate-api.com/docs",
    "terms_of_use": "https://www.exchangerate-api.com/terms",
    "supported_codes": [
        [
            "AED",
            "UAE Dirham"
        ],
        [
            "AFN",
            "Afghan Afghani"
        ],
        [
            "ALL",
            "Albanian Lek"
        ],
        [
            "AMD",
            "Armenian Dram"
        ],
        [
            "ANG",
            "Netherlands Antillian Guilder"
        ],
        [
            "AOA",
            "Angolan Kwanza"
        ],
        [
            "ARS",
            "Argentine Peso"
        ],
        [
            "AUD",
            "Australian Dollar"
        ],
        [
            "AWG",
            "Aruban Florin"
        ],
        [
            "AZN",
            "Azerbaijani Manat"
        ],
        [
            "BAM",
            "Bosnia and Herzegovina Convertible Mark"
        ],
        [
            "BBD",
            "Barbados Dollar"
        ],
        [
            "BDT",
            "Bangladeshi Taka"
        ],
        [
            "BGN",
            "Bulgarian Lev"
        ],
        [
            "BHD",
            "Bahraini Dinar"
        ],
        [
            "BIF",
            "Burundian Franc"
        ],
        [
            "BMD",
            "Bermudian Dollar"
        ],
        [
            "BND",
            "Brunei Dollar"
        ],
        [
            "BOB",
            "Bolivian Boliviano"
        ],
        [
            "BRL",
            "Brazilian Real"
        ],
        [
            "BSD",
            "Bahamian Dollar"
        ],
        [
            "BTN",
            "Bhutanese Ngultrum"
        ],
        [
            "BWP",
            "Botswana Pula"
        ],
        [
            "BYN",
            "Belarusian Ruble"
        ],
        [
            "BZD",
            "Belize Dollar"
        ],
        [
            "CAD",
            "Canadian Dollar"
        ],
        [
            "CDF",
            "Congolese Franc"
        ],
        [
            "CHF",
            "Swiss Franc"
        ],
        [
            "CLP",
            "Chilean Peso"
        ],
        [
            "CNY",
            "Chinese Renminbi"
        ],
        [
            "COP",
            "Colombian Peso"
        ],
        [
            "CRC",
            "Costa Rican Colon"
        ],
        [
            "CUP",
            "Cuban Peso"
        ],
        [
            "CVE",
            "Cape Verdean Escudo"
        ],
        [
            "CZK",
            "Czech Koruna"
        ],
        [
            "DJF",
            "Djiboutian Franc"
        ],
        [
            "DKK",
            "Danish Krone"
        ],
        [
            "DOP",
            "Dominican Peso"
        ],
        [
            "DZD",
            "Algerian Dinar"
        ],
        [
            "EGP",
            "Egyptian Pound"
        ],
        [
            "ERN",
            "Eritrean Nakfa"
        ],
        [
            "ETB",
            "Ethiopian Birr"
        ],
        [
            "EUR",
            "Euro"
        ],
        [
            "FJD",
            "Fiji Dollar"
        ],
        [
            "FKP",
            "Falkland Islands Pound"
        ],
        [
            "FOK",
            "Faroese Króna"
        ],
        [
            "GBP",
            "Pound Sterling"
        ],
        [
            "GEL",
            "Georgian Lari"
        ],
        [
            "GGP",
            "Guernsey Pound"
        ],
        [
            "GHS",
            "Ghanaian Cedi"
        ],
        [
            "GIP",
            "Gibraltar Pound"
        ],
        [
            "GMD",
            "Gambian Dalasi"
        ],
        [
            "GNF",
            "Guinean Franc"
        ],
        [
            "GTQ",
            "Guatemalan Quetzal"
        ],
        [
            "GYD",
            "Guyanese Dollar"
        ],
        [
            "HKD",
            "Hong Kong Dollar"
        ],
        [
            "HNL",
            "Honduran Lempira"
        ],
        [
            "HRK",
            "Croatian Kuna"
        ],
        [
            "HTG",
            "Haitian Gourde"
        ],
        [
            "HUF",
            "Hungarian Forint"
        ],
        [
            "IDR",
            "Indonesian Rupiah"
        ],
        [
            "ILS",
            "Israeli New Shekel"
        ],
        [
            "IMP",
            "Manx Pound"
        ],
        [
            "INR",
            "Indian Rupee"
        ],
        [
            "IQD",
            "Iraqi Dinar"
        ],
        [
            "IRR",
            "Iranian Rial"
        ],
        [
            "ISK",
            "Icelandic Króna"
        ],
        [
            "JEP",
            "Jersey Pound"
        ],
        [
            "JMD",
            "Jamaican Dollar"
        ],
        [
            "JOD",
            "Jordanian Dinar"
        ],
        [
            "JPY",
            "Japanese Yen"
        ],
        [
            "KES",
            "Kenyan Shilling"
        ],
        [
            "KGS",
            "Kyrgyzstani Som"
        ],
        [
            "KHR",
            "Cambodian Riel"
        ],
        [
            "KID",
            "Kiribati Dollar"
        ],
        [
            "KMF",
            "Comorian Franc"
        ],
        [
            "KRW",
            "South Korean Won"
        ],
        [
            "KWD",
            "Kuwaiti Dinar"
        ],
        [
            "KYD",
            "Cayman Islands Dollar"
        ],
        [
            "KZT",
            "Kazakhstani Tenge"
        ],
        [
            "LAK",
            "Lao Kip"
        ],
        [
            "LBP",
            "Lebanese Pound"
        ],
        [
            "LKR",
            "Sri Lanka Rupee"
        ],
        [
            "LRD",
            "Liberian Dollar"
        ],
        [
            "LSL",
            "Lesotho Loti"
        ],
        [
            "LYD",
            "Libyan Dinar"
        ],
        [
            "MAD",
            "Moroccan Dirham"
        ],
        [
            "MDL",
            "Moldovan Leu"
        ],
        [
            "MGA",
            "Malagasy Ariary"
        ],
        [
            "MKD",
            "Macedonian Denar"
        ],
        [
            "MMK",
            "Burmese Kyat"
        ],
        [
            "MNT",
            "Mongolian Tögrög"
        ],
        [
            "MOP",
            "Macanese Pataca"
        ],
        [
            "MRU",
            "Mauritanian Ouguiya"
        ],
        [
            "MUR",
            "Mauritian Rupee"
        ],
        [
            "MVR",
            "Maldivian Rufiyaa"
        ],
        [
            "MWK",
            "Malawian Kwacha"
        ],
        [
            "MXN",
            "Mexican Peso"
        ],
        [
            "MYR",
            "Malaysian Ringgit"
        ],
        [
            "MZN",
            "Mozambican Metical"
        ],
        [
            "NAD",
            "Namibian Dollar"
        ],
        [
            "NGN",
            "Nigerian Naira"
        ],
        [
            "NIO",
            "Nicaraguan Córdoba"
        ],
        [
            "NOK",
            "Norwegian Krone"
        ],
        [
            "NPR",
            "Nepalese Rupee"
        ],
        [
            "NZD",
            "New Zealand Dollar"
        ],
        [
            "OMR",
            "Omani Rial"
        ],
        [
            "PAB",
            "Panamanian Balboa"
        ],
        [
            "PEN",
            "Peruvian Sol"
        ],
        [
            "PGK",
            "Papua New Guinean Kina"
        ],
        [
            "PHP",
            "Philippine Peso"
        ],
        [
            "PKR",
            "Pakistani Rupee"
        ],
        [
            "PLN",
            "Polish Złoty"
        ],
        [
            "PYG",
            "Paraguayan Guaraní"
        ],
        [
            "QAR",
            "Qatari Riyal"
        ],
        [
            "RON",
            "Romanian Leu"
        ],
        [
            "RSD",
            "Serbian Dinar"
        ],
        [
            "RUB",
            "Russian Ruble"
        ],
        [
            "RWF",
            "Rwandan Franc"
        ],
        [
            "SAR",
            "Saudi Riyal"
        ],
        [
            "SBD",
            "Solomon Islands Dollar"
        ],
        [
            "SCR",
            "Seychellois Rupee"
        ],
        [
            "SDG",
            "Sudanese Pound"
        ],
        [
            "SEK",
            "Swedish Krona"
        ],
        [
            "SGD",
            "Singapore Dollar"
        ],
        [
            "SHP",
            "Saint Helena Pound"
        ],
        [
            "SLE",
            "Sierra Leonean Leone"
        ],
        [
            "SLL",
            "Sierra Leonean Leone"
        ],
        [
            "SOS",
            "Somali Shilling"
        ],
        [
            "SRD",
            "Surinamese Dollar"
        ],
        [
            "SSP",
            "South Sudanese Pound"
        ],
        [
            "STN",
            "São Tomé and Príncipe Dobra"
        ],
        [
            "SYP",
            "Syrian Pound"
        ],
        [
            "SZL",
            "Eswatini Lilangeni"
        ],
        [
            "THB",
            "Thai Baht"
        ],
        [
            "TJS",
            "Tajikistani Somoni"
        ],
        [
            "TMT",
            "Turkmenistan Manat"
        ],
        [
            "TND",
            "Tunisian Dinar"
        ],
        [
            "TOP",
            "Tongan Paʻanga"
        ],
        [
            "TRY",
            "Turkish Lira"
        ],
        [
            "TTD",
            "Trinidad and Tobago Dollar"
        ],
        [
            "TVD",
            "Tuvaluan Dollar"
        ],
        [
            "TWD",
            "New Taiwan Dollar"
        ],
        [
            "TZS",
            "Tanzanian Shilling"
        ],
        [
            "UAH",
            "Ukrainian Hryvnia"
        ],
        [
            "UGX",
            "Ugandan Shilling"
        ],
        [
            "USD",
            "United States Dollar"
        ],
        [
            "UYU",
            "Uruguayan Peso"
        ],
        [
            "UZS",
            "Uzbekistani So'm"
        ],
        [
            "VES",
            "Venezuelan Bolívar Soberano"
        ],
        [
            "VND",
            "Vietnamese Đồng"
        ],
        [
            "VUV",
            "Vanuatu Vatu"
        ],
        [
            "WST",
            "Samoan Tālā"
        ],
        [
            "XAF",
            "Central African CFA Franc"
        ],
        [
            "XCD",
            "East Caribbean Dollar"
        ],
        [
            "XDR",
            "Special Drawing Rights"
        ],
        [
            "XOF",
            "West African CFA franc"
        ],
        [
            "XPF",
            "CFP Franc"
        ],
        [
            "YER",
            "Yemeni Rial"
        ],
        [
            "ZAR",
            "South African Rand"
        ],
        [
            "ZMW",
            "Zambian Kwacha"
        ],
        [
            "ZWL",
            "Zimbabwean Dollar"
        ]
    ]
}
`

type Currency struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Response struct {
	Result         string     `json:"result"`
	SupportedCodes [][]string `json:"supported_codes"`
}

const version = "1.0.0"

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
}

//	@title			Exchanger API
//	@version		1.0
//	@description	Exchanger Open API
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/v1

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	cfg := config{
		port: env.GetInt("PORT", 8080),
		env:  env.GetString("ENV", "development"),
		dbConfig: dbConfig{
			dsn:         env.GetString("DATABASE_URL", "postgres://root:secret@localhost:5432/exchanger?sslmode=disable"),
			maxIdleConn: env.GetInt("MAX_IDLE_CONN", 25),
			maxOpenConn: env.GetInt("MAX_OPEN_CONN", 25),
			maxIdleTime: env.GetString("MAX_IDLE_TIME", "15m"),
		},
		mailConfig: mailConfig{
			sender:   env.GetString("MAIL_SENDER", "Exchanger"),
			host:     env.GetString("MAIL_HOST", ""),
			port:     env.GetInt("MAIL_PORT", 25),
			username: env.GetString("MAIL_USERNAME", ""),
			password: env.GetString("MAIL_PASSWORD", ""),
		},
		jwtConfig: jwtConfig{
			issuer:        env.GetString("JWT_ISSUER", "Exchanger"),
			secret:        env.GetString("JWT_SECRET", ""),
			expiry:        env.GetString("JWT_EXPIRY", "15"),
			refreshExpiry: env.GetString("JWT_REFRESH_EXPIRY", "72h")},
	}

	// Logger
	opts := PrettyHandlerOptions{
		options: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := NewPrettyHandler(os.Stdout, opts)
	logger := slog.New(handler)

	// Database
	db, err := database.ConnectDB(
		cfg.dbConfig.dsn,
		cfg.dbConfig.maxIdleConn,
		cfg.dbConfig.maxOpenConn,
		cfg.dbConfig.maxIdleTime,
	)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// Mailer
	mailer := mail.NewMailer(
		cfg.mailConfig.sender,
		cfg.mailConfig.host,
		cfg.mailConfig.port,
		cfg.mailConfig.username,
		cfg.mailConfig.password,
	)

	// Storage (repository)
	storage := store.NewStorage(db)

	app := application{
		config: cfg,
		store:  storage,
		mailer: mailer,
		logger: logger,
	}

	// Serve application
	log.Fatal(app.serve())
}
