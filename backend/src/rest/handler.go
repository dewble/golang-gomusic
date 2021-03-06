package rest

// 패키지 선언하고 외부 패키지 임포트
import (
	"fmt"
	"gomusic/backend/src/dblayer"
	"gomusic/backend/src/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
)

// 코드 확장성을 높이고자 핸들러의 모든 메서드를 포함하는 인터페이스를 만든다.
type HandlerInterface interface {
	GetProducts(c *gin.Context)
	GetPromos(c *gin.Context)
	AddUser(c *gin.Context)
	SignIn(c *gin.Context)
	SignOut(c *gin.Context)
	GetOrders(c *gin.Context)
	Charge(c *gin.Context)
}

// 모든 메서드가 있는 Handler 구조체 정의
// Handler 타입은 데이터를 읽거나 수정하기 때문에 데이터베이스 레이어 인터페이스에 접근할 수 있어야한다.
type Handler struct {
	db dblayer.DBLayer
}

// 좋은 설계 원칙에 따라 Handler 생성자를 만든다
// 데이터베이스 레이어 타입의 초기화를 위해 이 생성자의 구현을 앞으로 계속 추가한다
func NewHandler() (HandlerInterface, error) {
	db, err := dblayer.NewORM("mysql", "gomusic:gomusic123@/gomusic")
	if err != nil {
		return nil, err
	}
	// Handler 객체에 대한 포인터 생성
	return &Handler{
		db: db,
	}, nil
}

func NewHandlerWithDB(db dblayer.DBLayer) HandlerInterface {
	return &Handler{db: db}
}

func (h *Handler) GetMainPage(c *gin.Context) {
	log.Println("Main page....")
	c.String(http.StatusOK, "Main page for secure API!!")
	//fmt.Fprintf(c.Writer, "Main page for secure API!!")
}

// 상품 목록 조회
// *gin.Context 타입 인자를 전달받는 GetProducts 메서드를 정의
func (h *Handler) GetProducts(c *gin.Context) {

	// DB 인터페이스가 nil이 아닌 값으로 초기화 됐는지 확인.
	// 이 객체를 통해 상품 목록을 조회
	if h.db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server database error"})
		return
	}
	products, err := h.db.GetAllProducts()

	// 에러가 발생한다면 HTTP 상태 코드를 포함한 JSON 데이터 반환
	if err != nil {
		/*
			첫 번째 인자는 HTTP 상태코드, 두 번째는 응답의 바디
		*/
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 에러가 발생하지 않았다면 데이터베이스에서 읽은 상품 반환, 데이터 모델에 JSON구조체 태그로 정의한 필드는 JSON 형식에 맞춰 변환
	fmt.Printf("Found %d products\n", len(products))
	c.JSON(http.StatusOK, products)
}

// 프로모션 목록 조회
// *gin.Context 타입 인자를 전달받는 GetPromos 메서드를 정의
func (h *Handler) GetPromos(c *gin.Context) {

	// DB 인터페이스가 nil이 아닌 값으로 초기화 됐는지 확인.
	// 이 객체를 통해 상품 목록을 조회

	if h.db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server database error"})
		return
	}
	promos, err := h.db.GetPromos()

	// 에러가 발생한다면 HTTP 상태 코드를 포함한 JSON 데이터 반환
	if err != nil {
		/*
			첫 번째 인자는 HTTP 상태코드, 두 번째는 응답의 바디
		*/
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// 에러가 발생하지 않았다면 데이터베이스에서 읽은 상품 반환, 데이터 모델에 JSON구조체 태그로 정의한 필드는 JSON 형식에 맞춰 변환
	c.JSON(http.StatusOK, promos)
}

func (h *Handler) AddUser(c *gin.Context) {
	if h.db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server database error"})
		return
	}
	var customer models.Customer

	// HTTP 요청 바디에서 JSON 문서를 추출하고 객체로 디코딩한다.
	// 아래의 경우 이 객체는 고객 데이터 모델을 나타내는 *models.Customer 타입이다
	err := c.ShouldBindJSON(&customer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// JSON 문서를 데이터 모델로 디코딩하고 AddUser 데이터베이스 레이어 메서드를 호출하고 데이터베이스에 로그인 상태를 저장하거나 신규 사용자를 추가
	customer, err = h.db.AddUser(customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customer)
}

// 사용자 로그인과 신규 가입
func (h *Handler) SignIn(c *gin.Context) {
	if h.db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server database error"})
		return
	}
	var customer models.Customer

	// HTTP 요청 바디에서 JSON 문서를 추출하고 객체로 디코딩한다.
	// 아래의 경우 이 객체는 고객 데이터 모델을 나타내는 *models.Customer 타입이다
	err := c.ShouldBindJSON(&customer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// JSON 문서를 데이터 모델로 디코딩하고 SignInUser 데이터베이스 레이어 메서드를 호출하고 데이터베이스에 로그인 상태를 저장하거나 신규 사용자를 추가
	customer, err = h.db.SignInUser(customer.Email, customer.Pass)
	if err != nil {

		// 잘못된 패스워드인 경우 forbiiden http 에러 반환
		if err == dblayer.ErrINVALIDPASSWORD {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customer)
}

func (h *Handler) SignOut(c *gin.Context) {
	if h.db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server database error"})
		return
	}

	// URL에서 로그아웃하는 사용자의 ID를 추출한다. *gin.Context 타입의 Param() 메서드를 사용한다.
	p := c.Param("id")
	// p는 문자형. 저수형으로 변환
	id, err := strconv.Atoi(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// SignOutUserById 데이터베이스 레이어 메서드를 호출하고 데이터베이스에 해당 사용자를 로그아웃 상태로 설정한다.
	err = h.db.SignOutUserById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

// 사용자의 주문 내역 조회
func (h *Handler) GetOrders(c *gin.Context) {
	if h.db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server database error"})
		return
	}

	// URL에서 로그아웃하는 사용자의 ID를 추출한다. *gin.Context 타입의 Param() 메서드를 사용한다.
	// id 매개변수 추출
	p := c.Param("id")
	// p는 문자형. 저수형으로 변환
	id, err := strconv.Atoi(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// GetCustomerOrdersByID 데이터베이스 레이어 메서드를 호출하고 주문 내역 조회
	orders, err := h.db.GetCustomerOrdersByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// 신용카드 결제 요청
func (h *Handler) Charge(c *gin.Context) {
	if h.db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server database error"})
		return
	}

	// Go 구조체를 정의하는 동시에 초기화
	request := struct {
		models.Order
		Remember    bool   `json:"rememberCard"`
		UseExisting bool   `json:"useExisting"`
		Token       string `json:"token"`
	}{}

	err := c.ShouldBindJSON(&request)
	log.Printf("request: %+v \n", request)
	// 파싱 중 에러 발생 시 보고 후 반환
	if err != nil {
		// JSON 형식의 요청 데이터를 request 구조체로 변환
		c.JSON(http.StatusBadRequest, request)
		return
	}

	// Set your secret key: remember to change this to your live secret key in production
	// Keys can be obtained from: https://dashboard.stripe.com/account/apikeys
	// They key below is just for testing
	stripe.Key = "sk_test_51JqdjpHqNVgKzGGZR1CaBtIw3eOw4HvXZh4rvhynS3BSoiaIZ4GTnGXubXs6yjUb7bYfLuFIPhknozJ9aaOmDYRa00xcnUrXwz"
	//test cards available at:	https://stripe.com/docs/testing#cards
	//setting charge parameters

	chargeP := &stripe.ChargeParams{
		// 요청에 명시된 판매 가격
		Amount: stripe.Int64(int64(request.Price)),
		// 결제 통화
		Currency: stripe.String("usd"),
		// 설명
		Description: stripe.String("GoMusic charge..."),
	}

	// 스트라이프 사용자 ID 초기화
	stripeCustomerID := ""

	//Either remembercard or use exeisting should be enabled but not both
	if request.UseExisting {
		// 저장된 카드 사용
		log.Println("Getting credit card id...")
		// 스트라이프 사용자 ID를 데이터베이스에서 조회하는 메서드
		stripeCustomerID, err = h.db.GetCreditCardCID(request.CustomerID)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		cp := &stripe.CustomerParams{}
		cp.SetSource(request.Token)
		customer, err := customer.New(cp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		stripeCustomerID = customer.ID
		if request.Remember {
			// 스트라이프 사용자 id를 저장하고 데이터베이스에 저장된 사용자 ID와 연결한다.
			err = h.db.SaveCreditCardForCustomer(request.CustomerID, stripeCustomerID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	//we should check if the customer already ordered the same item or not but for simplicity, let's assume it's a new order
	// 동일 상품 주문 여부 확인 없이 새로운 주문으로 가정
	// *stripe.ChargeParams 타입 인스턴스에 스트라이프 사용자 ID를 설정한다.
	chargeP.Customer = stripe.String(stripeCustomerID)
	// 신용카드 결제 요청
	_, err = charge.New(chargeP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.db.AddOrder(request.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

}
