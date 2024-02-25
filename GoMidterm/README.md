# Apple stores project


## Stores REST API
```
POST /menus
GET /menus/:id
PUT /menus/:id
DELETE /menus/:id
```

## DB Structure

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
Table restaurants_and_menu {
  id bigserial [primary key]
  created_at timestamp
  updated_at timestamp
  store bigserial
  product bigserial
}

Ref: restaurants_and_menu.restaurant < restaurants.id
Ref: restaurants_and_menu.menu < menu.id

```
