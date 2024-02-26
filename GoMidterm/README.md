# Apple stores project


## Stores REST API
```
POST /products
GET /products/:id
PUT /products/:id
DELETE /products/:id
```

## DB Structure

![image](https://github.com/kim0111/Go/assets/86676168/7c976f5d-e98a-4e9e-b51e-97b533d1c949)


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
