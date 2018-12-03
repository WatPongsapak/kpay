# KPAY Quiz

create an api for merchant seller and selling report.

## Requirement
- language: golang
- database: mongodb

## Merchant Fields
- Name
- Bank Account 
- Username
- Password

## Product Fields
- Name
- Amount
- Amount Changes (as history when update amount)

## APIs
| Name                 | Method | Endpoint                          |
|----------------------|--------|-----------------------------------|
| Register Merchant    | POST   | /merchant/register                |
| Merchant Information | GET    | /merchant/:id                     |
| Update Merchant      | POST   | /merchant/:id                     |
| List All Products    | GET    | /merchant/:id/products            |
| Add Product          | POST   | /merchant/:id/product             |
| Update Product       | POST   | /merchant/:id/product/:product_id |
| Remove Product       | DELETE | /merchant/:id/product/:product_id |
| Sell Reports         | POST   | /merchant/:id/report              |
| Buy Product          | POST   | /buy/product                      |

### Register Merchant
- auto generate username and password
- each api must check username/password except register and buy product
- cannot register using the same bank account

### Merchant Information
- response merchant information

### Update Merchant
- can only update name

### List All Products
- list all merchant products with name and amount

### Add Product
- add product for each merchant 
- amount can be present in 2 precision, ex. 100.01, 250.35
- maximum products is 5

### Update Product
- can only update amount
- in case of user already brought product, in sell report must calculate by old amount

### Remove Product
- remove product by product id
- cannot remove if user already brought product

### Sell Reports
- sell report range only by date
- provide list of selling products and amount accumulate with precision point 2 digit
- ensure there is index in all related fields collections,  must prove that there is no table scan

```json
{
	"date": "2018-11-01",
	"products": [
		{"name": "ABC", "selling_volume": 10},
		{"name": "DEF", "selling_volume": 5}
	],
	"accumulate": 100.25 
}
```

### Buy Product
- buy product from merchant with volume
