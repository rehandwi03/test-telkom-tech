CREATE TABLE public.cart_products (
                                      cart_id uuid NOT NULL,
                                      product_id uuid NOT NULL,
                                      quantity integer NOT NULL
);


--
-- Name: carts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.carts (
                              id uuid NOT NULL,
                              full_name character varying(50) NOT NULL
);


--
-- Name: files; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.files (
                              id uuid NOT NULL,
                              name character varying(50) NOT NULL,
                              location text NOT NULL,
                              bucket_name character varying(50) NOT NULL
);


--
-- Name: products; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.products (
                                 id uuid NOT NULL,
                                 name character varying(50) NOT NULL,
                                 price numeric(21,2) NOT NULL,
                                 description text NOT NULL,
                                 is_discount boolean NOT NULL,
                                 start_date_discount date,
                                 end_date_discount date,
                                 discount_value numeric(21,2)
);