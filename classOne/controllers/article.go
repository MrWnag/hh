package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"classOne/models"
	"math"
)

type ArticleController struct {
	beego.Controller
}


//商品列表页
func (this*ArticleController)ShowArticleList(){
	//1.查询


		//qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles)//1.pagesize  一页显示多少   2.start起始位置
		o := orm.NewOrm()
		//获取类型数据
		var types []models.GoodsType
		o.QueryTable("GoodsType").All(&types)
		this.Data["types"] = types


		//根据类型获取数据
		//1.接受数据
		//1.有一个orm对象

		qs :=o.QueryTable("GoodsSKU")

		pageIndex,err := this.GetInt("pageIndex")
		if err !=nil{
			pageIndex = 1//设置默认访问首页
		}
		typeName:=this.GetString("typeName")
		//beego.Info(typeName)
		//2.处理数据
		pageSize := 2 //定义一页显示多少条数据
		start:=pageSize*(pageIndex -1)//获取查询数据的起始位置
		var count int64
		var goodswithtype []models.GoodsSKU
		if typeName == ""{
			count ,_= qs.RelatedSel("GoodsType").Count()//返回数据条目数   加过滤器


			qs.Limit(pageSize,start).RelatedSel("GoodsType").All(&goodswithtype)
		}else {
			count ,_= qs.RelatedSel("GoodsType").Filter("GoodsType__Name",typeName).Count()//返回数据条目数   加过滤器
			qs.Limit(pageSize,start).RelatedSel("GoodsType").Filter("GoodsType__Name",typeName).All(&goodswithtype)
		}
		pageCount := float64(count)/float64(pageSize)  //求总页数
		pageCount=math.Ceil(pageCount)  //天花板函数获取正确总页码


		FirstPage := false//标识是否是首页
		EndPage := false //标识是否是末页
		//首页末页数据处理
		if pageIndex == 1{
			FirstPage = true
		}
		if pageIndex == int(pageCount){
			EndPage = true
		}
		//3.查询数据

		userName:=this.GetSession("userName")
		this.Data["userName"] = userName


		this.Data["typeName"] = typeName
		this.Data["EndPage"] = EndPage
		this.Data["FirstPage"] = FirstPage
		this.Data["count"] = count
		this.Data["pageCount"] = pageCount
		this.Data["pageIndex"] = pageIndex
		this.Data["goods"] = goodswithtype


		//2.把数据传递给视图显示
		this.Layout = "layout.html"
		this.TplName = "index.html"
}

//展示添加商品界面
func(this*ArticleController)ShowAddArticle(){
	//查询类型数据，传递到视图中
	o:=orm.NewOrm()
	var types []models.GoodsType
	o.QueryTable("GoodsType").All(&types)
	this.Data["types"] = types

	this.TplName = "add.html"
}

//处理增加商品业务
/*
1.那数据
2.判断数据
3.插入数据
4.返回试图
 */
func (this*ArticleController)HandleAddArtcile(){
	//1.那数据
	//那标题
	goodsName:= this.GetString("goodsName")
	goodsDesc := this.GetString("desc")
	f,h,err:=this.GetFile("uploadname")

	defer f.Close()
	//上传文件处理
	//1.判断文件格式
	ext := path.Ext(h.Filename)
	if ext != ".jpg" && ext != ".png"&&ext != ".jpeg"{
		beego.Info("上传文件格式不正确")
		return
	}

	//2.文件大小
	if h.Size>5000000{
		beego.Info("文件太大，不允许上传")
		return
	}

	//3.不能重名
	fileName := time.Now().Format("2006-01-02 15:04:05")


	err2:=this.SaveToFile("uploadname","./static/img/"+fileName+ext)
	if err != nil{
		beego.Info("上传文件失败")
		return
	}

	if err != nil{
		beego.Info("上传文件失败",err2)
		return
	}




	//3.插入数据
	//1.获取orm对象
	o := orm.NewOrm()
	//2.创建一个插入对象
	goods := models.GoodsSKU{}
	//3.赋值
	goods.Name = goodsName
	goods.Desc = goodsDesc
	goods.Image = "/static/img/"+fileName+ext


	//4.返回试图
	//给article对象复制
	//获取到下拉框传递过来的类型数据
	typeName:=this.GetString("selectType")
	//类型判断
	if typeName == ""{
		beego.Info("下拉匡数据错误")
		return
	}
	//获取type对象
	var goodsType models.GoodsType
	goodsType.Name = typeName
	err=o.Read(&goodsType,"Name")
	if err != nil{
		beego.Info("获取类型错误")
		return
	}
	goods.GoodsType = &goodsType

	//获取goodsSPU
	spuName := this.GetString("selectGoodsSPU")
	var goodsSPU models.Goods
	goodsSPU.Name = spuName
	err = o.Read(&goodsSPU,"Name")
	if err != nil{
		beego.Info("插入数据失败")
		return
	}
	goods.Goods = &goodsSPU

	//4.插入
	_,err = o.Insert(&goods)
	if err != nil{
		beego.Info("插入数据失败")
		return
	}
	this.Redirect("/Article/ShowArticle",302)

}



//显示商品详情
func (this*ArticleController)ShowContent(){
	//1.获取Id
	id,err:=this.GetInt("id")
	if err !=nil{
		beego.Info("查询数据为空")
		return
	}
	beego.Info(id)
	//2.查询数据
	//1.获取orm对象
	o:=orm.NewOrm()
	//2.获取查询对象
	var goods models.GoodsSKU
	//3.查询
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("Id",id).One(&goods)
	if err != nil{
		beego.Info("查询数据为空")
		return
	}
	//3.传递数据给视图
	this.Data["goods"] = goods
	this.Layout = "layout.html"
	this.TplName = "content.html"
}


//1.URLchuanzhi
//2.执行delete操作


//删除商品
func (this*ArticleController)HandleDelete(){
	id,err:=this.GetInt("id")
	if err != nil{
		beego.Info("获取商品错误")
		return
	}
	//1.orm对象
	o := orm.NewOrm()

	//要有删除对象
	article := models.GoodsSKU{Id:id}

	//3.删除
	o.Delete(&article)

	this.Redirect("/Article/ShowArticle",302)
}

//显示更新页面
func (this*ArticleController)ShowUpdate(){
	//获取数据
	id,err:= this.GetInt("id")
	if err != nil{
		beego.Info("连接错误")
		return
	}
	//查询操作
	o := orm.NewOrm()
	//2.获取查询对象
	var goods models.GoodsSKU
	//3.查询
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("Id",id).One(&goods)
	if err != nil{
		beego.Info("查询数据为空")
		return
	}

	//把数据传递给视图
	this.Data["article"] = goods
	this.Layout = "layout.html"
	this.TplName = "update.html"


}

//处理更新数据
func(this*ArticleController)HandleUpdate(){
	//1.拿数据
	name:=this.GetString("goodsName")
	desc := this.GetString("Desc")
	id,err:=this.GetInt("id")
	if err != nil {
		beego.Info("传输连接错误")
		return
	}

	//问题一 id是不是没有传过来
	//2.判断数据
	if name == "" || desc == ""{
		beego.Info("更新数据失败")
		return
	}
	f,h,err:=this.GetFile("uploadname")
	if err != nil{
		beego.Info("上传文件失败")
		return
	}
	defer f.Close()
	//1.判断大小
	if h.Size > 500000{
		beego.Info("图片太大")
		return
	}
	//2.判断类型
	ext:=path.Ext(h.Filename)
	if ext != ".jpg"&&ext!=".png"&&ext!=".jpeg"{
		beego.Info("上传文件类型错误")
		return
	}
	//3.防止文件名重复
	filename:=time.Now().Format("2006-01-02-15:04:05")
	this.SaveToFile("uploadname","./static/img/"+filename+ext)


	//更新操作
	o:=orm.NewOrm()
	goods := models.GoodsSKU{Id:id}
	//读取操作
	err = o.Read(&goods)
	if err != nil{
		beego.Info("要更新的文章不存在")
		return
	}
	//更新
	goods.Name = name
	goods.Desc = desc
	goods.Image = "/static/img/"+filename+ext
	_,err=o.Update(&goods)
	if err != nil{
		beego.Info("更新失败")
		return
	}

	//跳转
	this.Redirect("/Article/ShowArticle",302)

}


//展示添加商品类型
func (this*ArticleController)ShowAddType(){
	//1.读取类型表，显示数据
	o := orm.NewOrm()
	var goodsTypes[]models.GoodsType
	//查询
	_,err:=o.QueryTable("GoodsType").All(&goodsTypes)
	if err != nil{
		beego.Info("查询类型错误")
	}
	this.Data["title"] = "<title>添加类型</title>"
	this.Data["goodsTypes"] = goodsTypes
	this.Layout = "layout.html"
	this.TplName = "addType.html"
}
//处理添加类型业务
func (this*ArticleController)HandleAddType(){
	//1.获取数据
	typename:=this.GetString("typeName")
	if typename == ""{
		beego.Info("添加类型数据为空")
		return
	}
	//获取两张图片
	this.GetFile("uploadLogo")
	this.GetFile("uploadImage")

	//3.执行插入操作
	o := orm.NewOrm()
	var goodsType models.GoodsType
	goodsType.Name = typename
	_,err:=o.Insert(&goodsType)
	if err != nil{
		beego.Info("插入失败")
		return
	}
	//4.展示视图？
	this.Redirect("/Article/AddArticleType",302)
}

//退出登陆
func (this*ArticleController)Logout(){
	//1.删除登陆状态
	this.DelSession("userName")
	//2.跳转登陆页面
	this.Redirect("/",302)
}
//删除商品类型
func (this*ArticleController)DeleteType(){
	//1.获取类型Id
	id,err:=this.GetInt("id")

	//2.都要进行数据判断
	if err != nil{
		beego.Info("删除连接错误")
		return
	}



	//3.删除操作
	o := orm.NewOrm()
	artiType := models.GoodsType{Id:id}
	o.Delete(&artiType)

	//4.返回视图
	this.Redirect("/Article/AddArticleType",302)
}















































