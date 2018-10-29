package models

import ("github.com/astaxie/beego/orm"
_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct{ //用户表
	Id 			int
	Name 		string`orm:"size(20)"` 		//用户名
	PassWord 	string`orm:"size(20)"`		//登陆密码
	Email 		string`orm:"size(50)"`   	//邮箱
	Active 		bool`orm:"default(false)"`    	//是否激活  0 表示未激活  1表示激活
	Power 		int`orm:"default(0)"`		//权限设置
	Address		[]*Address `orm:"reverse(many)"`
	OrderInfo   []*OrderInfo `orm:"reverse(many)"`

}

type Address struct { //地址表
	Id int
	Receiver string `orm:"size(20)"`  			//收件人
	Addr string		`orm:"size(50)"`  			//收件地址
	Zip_code string  `orm:"size(20)"` 			//邮编
	Phone string	 `orm:"size(20)"` 			//联系方式
	Is_default bool	 `orm:"defalt(false)"`			//是否默认 0 为非默认  1为默认
	User *User `orm:"rel(fk)"` 	//用户ID
	OrderInfo   []*OrderInfo `orm:"reverse(many)"`
}

type Goods struct { //商品SPU表
	Id 		int
	Name 	string`orm:"size(20)"`  //商品名称
	Detail 	string`orm:"size(200)"` //详细描述
	GoodsSKU []*GoodsSKU `orm:"reverse(many)"`
}

type GoodsType struct{//商品类型表
	Id int
	Name string			//种类名称
	Logo string			//logo
	Image string   		//图片
	GoodsSKU []*GoodsSKU `orm:"reverse(many)"`
	IndexTypeGoodsBanner  []*IndexTypeGoodsBanner  `orm:"reverse(many)"`
}

type GoodsSKU struct { //商品SKU表
	Id int
	GoodsType *GoodsType `orm:"rel(fk)"`  //商品所属种类
	Goods     *Goods 	 `orm:"rel(fk)"` //商品SPU
	Name       string					 //商品名称
	Desc       string					 //商品简介
	Price      int						 //商品价格
	Unite      string					 //商品单位
	Image      string				 	 //商品图片
	Stock      int	`orm:"default(1)"`	 //商品库存
	Sales      int	`orm:"default(0)"`	 //商品销量
	Status     int	 `orm:"default(1)"`	 //商品状态
	Time       time.Time `orm:"auto_now_add"`
	GoodsImage []*GoodsImage `orm:"reverse(many)"`
	IndexGoodsBanner   []*IndexGoodsBanner `orm:"reverse(many)"`
	IndexTypeGoodsBanner []*IndexTypeGoodsBanner  `orm:"reverse(many)"`
	OrderGoods   []*OrderGoods `orm:"reverse(many)"`
}

type GoodsImage struct { //商品图片表
	Id 			int
	Image 		string					//商品图片
	GoodsSKU 	*GoodsSKU   `orm:"rel(fk)"` //商品SKU
}
type IndexGoodsBanner struct { //首页轮播商品展示表
	Id 		  int
	GoodsSKU  *GoodsSKU	`orm:"rel(fk)"`	//商品sku
	Image     string					//商品图片
	Index     int  `orm:"default(0)"`   //展示顺序
}

type IndexTypeGoodsBanner struct {//首页分类商品展示表
	Id 				int
	GoodsType 		*GoodsType 	`orm:"rel(fk)"`			//商品类型
	GoodsSKU  		*GoodsSKU  	`orm:"rel(fk)"`			//商品sku
	Display_Type 	int   		`orm:"default(1)"`		//展示类型 0代表标题，1代表文字
	Index 			int   		`orm:"default(0)"`		//展示顺序
}

type IndexPromotionBanner struct {//首页促销商品展示表
	Id 		int
	Name 	string	`orm:"size(20)"`				//活动名称
	Url 	string	`orm:"size(50)"`				//活动链接
	Image 	string	`orm:"size(50)"`			//活动图片
	Index 	int  `orm:"default(0)"` //展示顺序
}

type OrderInfo struct {//订单表
	Id 				int
	User 			*User	`orm:"rel(fk)"`		//用户
	Address 		*Address`orm:"rel(fk)"`		//地址
	Pay_Method 		int							//付款方式
	Total_Count 	int		`orm:"default(0)"`	//商品数量
	Total_Price 	int							//商品总价
	Transit_Price 	int							//运费
	Order_status 	int							//订单状态
	Trade_No 		string						//支付编号

	OrderGoods   []*OrderGoods `orm:"reverse(many)"`
}

type OrderGoods struct {//订单商品表
	Id 			int
	OrderInfo 	*OrderInfo	`orm:"rel(fk)"`	//订单
	GoodsSKU 	[]*GoodsSKU	`orm:"rel(m2m)"`	//商品
	Count 		int		`orm:"default(1)"`	//商品数量
	Price 		int							//商品价格
	Comment 	string						//评论
}
func init() {
	// set default database
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/dalyfresh?charset=utf8")

	// register model
	orm.RegisterModel(new(User),new(Address),new(OrderGoods),new(OrderInfo),new(IndexPromotionBanner),new(IndexTypeGoodsBanner),new(IndexGoodsBanner),new(GoodsImage),new(GoodsSKU),new(GoodsType),new(Goods))

	// create table
	orm.RunSyncdb("default", false, true)
}
