create table if not exists vendors (
    id serial primary key,
    name text not null,
    telephone text not null
);

create table if not exists products (
    id serial primary key,
    name text not null,
    description text,
    price numeric not null,
    vendor_id integer not null references vendors(id)
);

create table if not exists vendees (
    id serial primary key,
    name text not null,
    telephone text not null
);

DO $$ BEGIN
    create type order_status  as enum  ('pending', 'cancelled', 'completed') ;
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

create table if not exists orders (
    id serial primary key,
    vendee_id integer not null references vendees(id),
    order_date date not null default current_date,
    total numeric not null,
    status order_status not null default 'pending'
);

create index if not exists idx_orders_vendee_id on orders(vendee_id);

create table if not exists order_products (
    order_id integer not null references orders(id),
    product_id integer not null references products(id),
    quantity integer not null,
    primary key (order_id, product_id)
);

create index if not exists idx_order_products_product_id on order_products(product_id);
create index if not exists idx_order_products_order_id on order_products(order_id);
