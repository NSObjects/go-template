/*
 * Generated from OpenAPI3 document
 * Module: Order
 */

package param





// OrderListOrdersRequest 获取订单列表
type OrderListOrdersRequest struct {

	Page int `json:"page" validate:"min=1"` // 页码

	Size int `json:"size" validate:"min=1,max=100"` // 每页数量


}



// OrderCreateOrderRequest 创建订单
type OrderCreateOrderRequest struct {



	Userid int `json:"userId" validate:"required"` // 

	Productid int `json:"productId" validate:"required"` // 

	Quantity int `json:"quantity" validate:"required" validate:"min=1"` // 


}




// OrderListItem 
type OrderListItem struct {

	Id int `json:"id"` // 

	Userid int `json:"userId"` // 

	Productid int `json:"productId"` // 

	Quantity int `json:"quantity"` // 

}

// OrderData 
type OrderData struct {

	Id int `json:"id"` // 

	Userid int `json:"userId"` // 

	Productid int `json:"productId"` // 

	Quantity int `json:"quantity"` // 

}

