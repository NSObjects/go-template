/*
 * Generated from OpenAPI3 document
 * Module: Product
 */

package param





// ProductListProductsRequest 获取产品列表
type ProductListProductsRequest struct {

	Page int `json:"page" validate:"min=1"` // 页码

	Size int `json:"size" validate:"min=1,max=100"` // 每页数量


}



// ProductCreateProductRequest 创建产品
type ProductCreateProductRequest struct {



	Name string `json:"name" validate:"required" validate:"min=1,max=100"` // 

	Price float64 `json:"price" validate:"required" validate:"min=0"` // 


}




// ProductListItem 
type ProductListItem struct {

	Id int `json:"id"` // 

	Name string `json:"name"` // 

	Price float64 `json:"price"` // 

}

// ProductData 
type ProductData struct {

	Id int `json:"id"` // 

	Name string `json:"name"` // 

	Price float64 `json:"price"` // 

}

