# Apple stores project


## Stores REST API
```
POST /products
GET /products/:id
PUT /products/:id
DELETE /products/:id
```

## DB Structure

![image](https://github.com/kim0111/Go/assets/86676168/fd062bbc-8dea-49fe-bafb-83e2c2ab47b1)



```
Table stores {
  id bigserial [primary key]
  created_at timestamp
  updated_at timestamp
  title text
  coordinates text
  address text
  number_of_branches text
}

Table products {
  id bigserial [primary key]
  created_at timestamp
  updated_at timestamp
  title text
  description text
  for_what_country text
  price text

}

// many-to-many
Table stores_and_products {
  id bigserial [primary key]
  created_at timestamp
  updated_at timestamp
  store bigserial
  product bigserial
}

Ref: stores_and_products.store < stores.id
Ref: stores_and_products.product < product.id

```

## Some info
```
1)Create user: http://localhost:8081/api/v1/users
{
    "name": "Jona",
    "email": "jona@gmail.com",
    "password": "qwerty123"
}

2)User activation: http://localhost:8081/api/v1/users/activated
{
    "token": ""
}

3)User login: http://localhost:8081/api/v1/users/login
{
    "email": "jona@gmail.com",
    "password": "qwerty123"
}

4)Delete product with permission: http://localhost:8081/api/v1/products/{id}
key: Authorization
value: Bearer {auth token}
```
