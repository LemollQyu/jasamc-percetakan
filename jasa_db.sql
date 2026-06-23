--
-- PostgreSQL database dump
--

\restrict EiM8dBeMW4rh5hRrcCbqBJQtAXpO8nHq8udUQQR2nURXwo9BtxzI3t1nv3cvVny

-- Dumped from database version 18.1
-- Dumped by pg_dump version 18.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: service_categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.service_categories (
    id bigint NOT NULL,
    name character varying(200) NOT NULL,
    description text NOT NULL,
    slug character varying(200) NOT NULL,
    is_active boolean DEFAULT true,
    meta jsonb,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.service_categories OWNER TO postgres;

--
-- Name: service_categories_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.service_categories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.service_categories_id_seq OWNER TO postgres;

--
-- Name: service_categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.service_categories_id_seq OWNED BY public.service_categories.id;


--
-- Name: service_media; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.service_media (
    id bigint NOT NULL,
    service_id bigint NOT NULL,
    url text NOT NULL,
    type character varying(50) NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.service_media OWNER TO postgres;

--
-- Name: service_media_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.service_media_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.service_media_id_seq OWNER TO postgres;

--
-- Name: service_media_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.service_media_id_seq OWNED BY public.service_media.id;


--
-- Name: service_spesification_values; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.service_spesification_values (
    id bigint CONSTRAINT service_specification_values_id_not_null NOT NULL,
    service_id bigint CONSTRAINT service_specification_values_service_id_not_null NOT NULL,
    spesification_id bigint CONSTRAINT service_specification_values_specification_id_not_null NOT NULL,
    value character varying(200) CONSTRAINT service_specification_values_value_not_null NOT NULL,
    additional_price numeric(12,2) DEFAULT 0,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.service_spesification_values OWNER TO postgres;

--
-- Name: service_specification_values_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.service_specification_values_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.service_specification_values_id_seq OWNER TO postgres;

--
-- Name: service_specification_values_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.service_specification_values_id_seq OWNED BY public.service_spesification_values.id;


--
-- Name: service_spesifications; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.service_spesifications (
    id bigint CONSTRAINT service_specifications_id_not_null NOT NULL,
    service_id bigint CONSTRAINT service_specifications_service_id_not_null NOT NULL,
    name character varying(200) CONSTRAINT service_specifications_name_not_null NOT NULL,
    input_type character varying(50) CONSTRAINT service_specifications_input_type_not_null NOT NULL,
    options jsonb,
    is_required boolean DEFAULT true,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.service_spesifications OWNER TO postgres;

--
-- Name: service_specifications_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.service_specifications_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.service_specifications_id_seq OWNER TO postgres;

--
-- Name: service_specifications_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.service_specifications_id_seq OWNED BY public.service_spesifications.id;


--
-- Name: services; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.services (
    id bigint NOT NULL,
    category_id bigint NOT NULL,
    name character varying(200) NOT NULL,
    slug character varying(200) NOT NULL,
    description text NOT NULL,
    base_price numeric(12,2),
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    duration_per_unit integer DEFAULT 0
);


ALTER TABLE public.services OWNER TO postgres;

--
-- Name: services_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.services_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.services_id_seq OWNER TO postgres;

--
-- Name: services_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.services_id_seq OWNED BY public.services.id;


--
-- Name: service_categories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_categories ALTER COLUMN id SET DEFAULT nextval('public.service_categories_id_seq'::regclass);


--
-- Name: service_media id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_media ALTER COLUMN id SET DEFAULT nextval('public.service_media_id_seq'::regclass);


--
-- Name: service_spesification_values id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_spesification_values ALTER COLUMN id SET DEFAULT nextval('public.service_specification_values_id_seq'::regclass);


--
-- Name: service_spesifications id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_spesifications ALTER COLUMN id SET DEFAULT nextval('public.service_specifications_id_seq'::regclass);


--
-- Name: services id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.services ALTER COLUMN id SET DEFAULT nextval('public.services_id_seq'::regclass);


--
-- Data for Name: service_categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.service_categories (id, name, description, slug, is_active, meta, created_at, updated_at) FROM stdin;
2	Spanduk & Banner	Kategori untuk layanan spanduk dan banner	spanduk-banner	t	{"icon": "localhost:8081/static/jasa/icon-category/1779185299519072100.webp"}	2026-02-17 00:00:13.022444+07	2026-02-17 00:00:13.022444+07
6	Undangan	Siap, ini beberapa versi deskripsi singkat buat undangan percetakan:\n\n1. Undangan berkualitas dengan desain elegan, cocok untuk berbagai acara spesial Anda.\n\n2. Cetak undangan premium dengan hasil tajam dan bahan terbaik.\n\n3. Undangan custom sesuai keinginan, tampil eksklusif dan berkesan.\n\n4. Solusi cetak undangan cepat, rapi, dan harga terjangkau.\n\n5. Buat momen spesial lebih berkesan dengan undangan berkualitas.\n\nKalau mau lebih spesifik (nikahan, ulang tahun, dll), bilang aja nanti gue sesuaikan biar lebih ngena \n	undangan	t	{"icon": "localhost:8081/static/jasa/icon-category/1779185571096725900.webp"}	2026-05-01 10:26:59.939253+07	2026-05-01 10:26:59.939253+07
5	Cetak Sticker	Cetak sticker berkualitas tinggi dengan hasil tajam dan warna yang tahan lama. Cocok untuk kebutuhan promosi, branding produk, label kemasan, maupun keperluan pribadi.	cetak-sticker	t	{"icon": "localhost:8081/static/jasa/icon-category/1779185627397720000.webp"}	2026-03-31 17:19:30.572633+07	2026-03-31 17:19:30.572633+07
1	Print On Paper	Layanan penyediaan media kertas untuk kebutuhan operasional, administrasi, dan dokumentasi. Mencakup cetak dokumen perkantoran, formulir fungsional, hingga materi edukasi yang mengutamakan ketajaman informasi dan standarisasi format dokumen.	print-on-paper	t	{"icon": "localhost:8081/static/jasa/icon-category/1779185258581360300.webp"}	2026-01-09 20:24:40.653688+07	2026-01-09 20:24:40.653688+07
3	Photograph	Kategori untuk semua layanan photo atau cetak foto di macam macam ukuran	photograph	t	{"icon": "localhost:8081/static/jasa/icon-category/1778257018244790800.webp"}	2026-02-17 18:02:29.788271+07	2026-02-17 18:02:29.788271+07
\.


--
-- Data for Name: service_media; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.service_media (id, service_id, url, type, created_at, updated_at) FROM stdin;
26	3	localhost:8081/static/jasa/icon-service/1770384516565833100.svg	icon	2026-02-06 20:28:36.568119+07	2026-02-06 20:28:36.568119+07
35	5	localhost:8081/static/jasa/icon-service/1771329662887427100.png	icon	2026-02-17 19:01:02.956559+07	2026-02-17 19:01:02.956559+07
43	6	localhost:8081/static/jasa/icon-service/1774952705596433600.png	icon	2026-03-31 17:25:05.614955+07	2026-03-31 17:25:05.614955+07
44	6	localhost:8081/static/jasa/thumbnail-service/1774952705598736600.png	thumbnail	2026-03-31 17:25:05.614955+07	2026-03-31 17:25:05.614955+07
45	6	localhost:8081/static/jasa/thumbnail-service/1774952705600340200.jpg	thumbnail	2026-03-31 17:25:05.614955+07	2026-03-31 17:25:05.614955+07
46	6	localhost:8081/static/jasa/gallery-service/1774952705601898600.jpg	gallery	2026-03-31 17:25:05.614955+07	2026-03-31 17:25:05.614955+07
47	7	localhost:8081/static/jasa/icon-service/1777607373068069100.png	icon	2026-05-01 10:49:33.149608+07	2026-05-01 10:49:33.149608+07
48	7	localhost:8081/static/jasa/thumbnail-service/1777607373070702100.jpg	thumbnail	2026-05-01 10:49:33.149608+07	2026-05-01 10:49:33.149608+07
49	7	localhost:8081/static/jasa/gallery-service/1777607373072276100.jpg	gallery	2026-05-01 10:49:33.149608+07	2026-05-01 10:49:33.149608+07
50	7	localhost:8081/static/jasa/gallery-service/1777607373074395000.jpg	gallery	2026-05-01 10:49:33.149608+07	2026-05-01 10:49:33.149608+07
51	7	localhost:8081/static/jasa/gallery-service/1777607373076004900.jpg	gallery	2026-05-01 10:49:33.149608+07	2026-05-01 10:49:33.149608+07
52	8	localhost:8081/static/jasa/icon-service/1777621122462454400.png	icon	2026-05-01 14:38:42.48212+07	2026-05-01 14:38:42.48212+07
53	8	localhost:8081/static/jasa/thumbnail-service/1777621122467422800.png	thumbnail	2026-05-01 14:38:42.48212+07	2026-05-01 14:38:42.48212+07
54	8	localhost:8081/static/jasa/gallery-service/1777621122470314000.png	gallery	2026-05-01 14:38:42.48212+07	2026-05-01 14:38:42.48212+07
55	8	localhost:8081/static/jasa/gallery-service/1777621122473208900.png	gallery	2026-05-01 14:38:42.48212+07	2026-05-01 14:38:42.48212+07
56	5	localhost:8081/static/jasa/thumbnail-service/1779195013111585600.webp	thumbnail	2026-05-19 19:50:13.114412+07	2026-05-19 19:50:13.114412+07
57	5	localhost:8081/static/jasa/thumbnail-service/1779195028067499800.webp	thumbnail	2026-05-19 19:50:28.069165+07	2026-05-19 19:50:28.069165+07
58	5	localhost:8081/static/jasa/thumbnail-service/1779195044232224200.jpg	thumbnail	2026-05-19 19:50:44.233276+07	2026-05-19 19:50:44.233276+07
59	5	localhost:8081/static/jasa/gallery-service/1779195095969200100.webp	gallery	2026-05-19 19:51:35.971379+07	2026-05-19 19:51:35.971379+07
60	5	localhost:8081/static/jasa/gallery-service/1779195122602665200.webp	gallery	2026-05-19 19:52:02.604243+07	2026-05-19 19:52:02.604243+07
61	1	localhost:8081/static/jasa/thumbnail-service/1779197474041082600.jpg	thumbnail	2026-05-19 20:31:14.043269+07	2026-05-19 20:31:14.043269+07
62	1	localhost:8081/static/jasa/icon-service/1779197546691119700.webp	icon	2026-05-19 20:32:26.693217+07	2026-05-19 20:32:26.693217+07
63	1	localhost:8081/static/jasa/gallery-service/1779197987077426200.jpg	gallery	2026-05-19 20:39:47.078449+07	2026-05-19 20:39:47.078449+07
64	1	localhost:8081/static/jasa/gallery-service/1779197993154017800.jpg	gallery	2026-05-19 20:39:53.155587+07	2026-05-19 20:39:53.155587+07
65	1	localhost:8081/static/jasa/thumbnail-service/1779198000652264500.jpg	thumbnail	2026-05-19 20:40:00.653871+07	2026-05-19 20:40:00.653871+07
66	3	localhost:8081/static/jasa/thumbnail-service/1779200103784855200.jpg	thumbnail	2026-05-19 21:15:03.786426+07	2026-05-19 21:15:03.786426+07
68	3	localhost:8081/static/jasa/thumbnail-service/1779200159979695500.jpg	thumbnail	2026-05-19 21:15:59.981442+07	2026-05-19 21:15:59.981442+07
69	3	localhost:8081/static/jasa/gallery-service/1779200170034372900.jpg	gallery	2026-05-19 21:16:10.036687+07	2026-05-19 21:16:10.036687+07
70	3	localhost:8081/static/jasa/gallery-service/1779200242367708300.jpg	gallery	2026-05-19 21:17:22.371494+07	2026-05-19 21:17:22.371494+07
71	3	localhost:8081/static/jasa/gallery-service/1779200254834928300.jpg	gallery	2026-05-19 21:17:34.841224+07	2026-05-19 21:17:34.841224+07
72	4	localhost:8081/static/jasa/gallery-service/1779200686458699800.png	gallery	2026-05-19 21:24:46.463107+07	2026-05-19 21:24:46.463107+07
73	4	localhost:8081/static/jasa/thumbnail-service/1779200904523651500.jpg	thumbnail	2026-05-19 21:28:24.613057+07	2026-05-19 21:28:24.613057+07
74	4	localhost:8081/static/jasa/icon-service/1779200957172366900.webp	icon	2026-05-19 21:29:17.174495+07	2026-05-19 21:29:17.174495+07
75	2	localhost:8081/static/jasa/thumbnail-service/1779201158206910200.PNG	thumbnail	2026-05-19 21:32:38.208467+07	2026-05-19 21:32:38.208467+07
76	2	localhost:8081/static/jasa/thumbnail-service/1779201194207109200.jpg	thumbnail	2026-05-19 21:33:14.208695+07	2026-05-19 21:33:14.208695+07
77	2	localhost:8081/static/jasa/gallery-service/1779201229458371000.jpg	gallery	2026-05-19 21:33:49.459372+07	2026-05-19 21:33:49.459372+07
78	2	localhost:8081/static/jasa/gallery-service/1779201239269294300.jpg	gallery	2026-05-19 21:33:59.271349+07	2026-05-19 21:33:59.271349+07
79	2	localhost:8081/static/jasa/icon-service/1779201443615362300.webp	icon	2026-05-19 21:37:23.617461+07	2026-05-19 21:37:23.617461+07
80	9	localhost:8081/static/jasa/icon-service/1779348399465181200.png	icon	2026-05-21 14:26:39.496428+07	2026-05-21 14:26:39.496428+07
81	9	localhost:8081/static/jasa/thumbnail-service/1779348399468314200.webp	thumbnail	2026-05-21 14:26:39.496428+07	2026-05-21 14:26:39.496428+07
82	9	localhost:8081/static/jasa/gallery-service/1779348399469528600.webp	gallery	2026-05-21 14:26:39.496428+07	2026-05-21 14:26:39.496428+07
83	10	localhost:8081/static/jasa/icon-service/1779961733480474400.webp	icon	2026-05-28 16:48:53.512895+07	2026-05-28 16:48:53.512895+07
84	10	localhost:8081/static/jasa/thumbnail-service/1779961733485131200.jpg	thumbnail	2026-05-28 16:48:53.512895+07	2026-05-28 16:48:53.512895+07
85	10	localhost:8081/static/jasa/thumbnail-service/1779961733487171700.jpg	thumbnail	2026-05-28 16:48:53.512895+07	2026-05-28 16:48:53.512895+07
86	10	localhost:8081/static/jasa/gallery-service/1779961733489209800.jpg	gallery	2026-05-28 16:48:53.512895+07	2026-05-28 16:48:53.512895+07
87	10	localhost:8081/static/jasa/gallery-service/1779961733491259800.jpg	gallery	2026-05-28 16:48:53.512895+07	2026-05-28 16:48:53.512895+07
88	10	localhost:8081/static/jasa/gallery-service/1779961733493821300.webp	gallery	2026-05-28 16:48:53.512895+07	2026-05-28 16:48:53.512895+07
\.


--
-- Data for Name: service_spesification_values; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.service_spesification_values (id, service_id, spesification_id, value, additional_price, created_at, updated_at) FROM stdin;
1	1	1	Sidu	100.00	2026-01-10 14:02:02.822907+07	2026-01-10 14:02:02.822907+07
3	1	1	Natural	0.00	2026-01-10 14:03:21.454817+07	2026-01-10 14:03:21.454817+07
2	1	2	Hitam-putih	0.00	2026-01-10 14:02:21.439537+07	2026-01-10 14:05:01.436197+07
4	1	2	Warna	500.00	2026-01-10 14:07:45.263831+07	2026-01-10 14:07:45.263831+07
5	1	3	A4	0.00	2026-01-10 14:08:59.835371+07	2026-01-10 14:08:59.835371+07
6	1	3	F4	0.00	2026-01-10 14:09:06.061475+07	2026-01-10 14:09:06.061475+07
7	2	6	Warna	500.00	2026-01-12 11:21:05.176257+07	2026-01-12 11:21:05.176257+07
8	2	6	Hitam-putih	0.00	2026-01-12 11:21:16.175485+07	2026-01-12 11:21:16.175485+07
9	2	7	A4	0.00	2026-01-12 11:21:43.646665+07	2026-01-12 11:21:43.646665+07
10	2	7	F4	0.00	2026-01-12 11:21:48.772264+07	2026-01-12 11:21:48.772264+07
11	2	8	Sidu	100.00	2026-01-12 11:22:53.102865+07	2026-01-12 11:22:53.102865+07
12	2	8	Natural	0.00	2026-01-12 11:23:01.889619+07	2026-01-12 11:23:01.889619+07
13	2	9	Sidu	100.00	2026-01-12 11:29:12.970045+07	2026-01-12 11:29:12.970045+07
14	2	9	Natural	0.00	2026-01-12 11:29:22.722537+07	2026-01-12 11:29:22.722537+07
15	2	10	A4	0.00	2026-01-12 11:29:44.685214+07	2026-01-12 11:29:44.685214+07
16	2	10	F4	0.00	2026-01-12 11:29:48.55572+07	2026-01-12 11:29:48.55572+07
17	2	11	Warna	500.00	2026-01-12 11:30:37.584171+07	2026-01-12 11:30:37.584171+07
18	2	11	Hitam-putih	0.00	2026-01-12 11:30:53.890019+07	2026-01-12 11:30:53.890019+07
19	2	12		1000.00	2026-01-12 11:34:01.70653+07	2026-01-12 11:34:01.70653+07
20	2	13		1000.00	2026-01-12 11:35:13.326384+07	2026-01-12 11:35:13.326384+07
21	4	14	Warna	500.00	2026-01-12 12:13:56.872278+07	2026-01-12 12:13:56.872278+07
22	4	14	Hitam-putih	0.00	2026-01-12 12:14:05.695004+07	2026-01-12 12:14:05.695004+07
23	4	15	A4	0.00	2026-01-12 12:14:26.704201+07	2026-01-12 12:14:26.704201+07
24	4	15	F4	0.00	2026-01-12 12:14:33.060984+07	2026-01-12 12:14:33.060984+07
25	4	16	Sidu	100.00	2026-01-12 12:15:10.866342+07	2026-01-12 12:15:10.866342+07
26	4	16	Natural	0.00	2026-01-12 12:15:24.35225+07	2026-01-12 12:15:24.35225+07
27	5	18	2x3	0.00	2026-02-19 14:34:19.326804+07	2026-02-19 14:34:19.326804+07
28	5	18	3x4	250.00	2026-02-19 14:36:04.61453+07	2026-02-19 14:36:04.61453+07
29	5	18	4x6	500.00	2026-02-19 14:36:50.335678+07	2026-02-19 14:36:50.335678+07
30	6	20	A4	2000.00	2026-03-31 17:27:12.916847+07	2026-03-31 17:27:12.916847+07
31	6	20	F4	3000.00	2026-03-31 17:27:23.111163+07	2026-03-31 17:27:23.111163+07
32	6	21	Kertas	0.00	2026-03-31 17:28:50.304957+07	2026-03-31 17:28:50.304957+07
33	6	23	Ya	1000.00	2026-03-31 17:30:28.941398+07	2026-03-31 17:30:28.941398+07
34	7	24	Undangan Fadhil RS	0.00	2026-05-01 11:07:25.791066+07	2026-05-01 11:07:25.791066+07
35	7	25	Tanpa Perlengkapan	0.00	2026-05-01 11:10:45.587539+07	2026-05-01 11:10:45.587539+07
36	7	25	Label Nama	200.00	2026-05-01 11:12:28.597036+07	2026-05-01 11:12:28.597036+07
37	7	25	Plastik Undangan	300.00	2026-05-01 11:14:29.901618+07	2026-05-01 11:14:29.901618+07
38	7	25	Label Nama dan Plastik	500.00	2026-05-01 11:14:47.468736+07	2026-05-01 11:14:47.468736+07
39	8	26	Undangan Maliq 98	0.00	2026-05-01 15:41:35.483629+07	2026-05-01 15:41:35.483629+07
40	8	27	Tanpa Perlengkapan	0.00	2026-05-05 12:43:58.184595+07	2026-05-05 12:43:58.184595+07
41	8	27	Label Nama	200.00	2026-05-05 12:44:07.331431+07	2026-05-05 12:44:07.331431+07
42	8	27	Plastik Undangan	300.00	2026-05-05 12:44:14.223657+07	2026-05-05 12:44:14.223657+07
43	8	27	Label Nama dan Plastik	500.00	2026-05-05 12:44:21.240383+07	2026-05-05 12:44:21.240383+07
44	9	28	mini banner 25 x 40 - Art Carton	0.00	2026-05-21 14:30:59.817926+07	2026-05-21 14:30:59.817926+07
45	9	28	mini banner 25 x 40 - Art Carton Glossy	5000.00	2026-05-21 14:31:10.069896+07	2026-05-21 14:31:10.069896+07
46	9	28	mini banner 25 x 40 - Art Carton Doft	5000.00	2026-05-21 14:31:17.324981+07	2026-05-21 14:31:17.324981+07
48	5	29	sepaket 3x4 isi 4	4000.00	2026-05-23 12:24:04.672664+07	2026-05-23 12:24:04.672664+07
47	5	29	sepaket 2x3 isi 6	4000.00	2026-05-23 12:23:57.636815+07	2026-05-23 12:24:25.449648+07
49	5	29	sepaket 4x6 isi 4	5000.00	2026-05-23 12:24:34.127976+07	2026-05-23 12:24:34.127976+07
50	10	30	Flexi China	22000.00	2026-05-28 16:51:51.67085+07	2026-05-28 16:52:59.452188+07
51	10	30	Flexi Korea	28000.00	2026-05-28 16:52:05.26354+07	2026-05-28 16:53:14.945484+07
52	10	30	Flexi Jerman	38000.00	2026-05-28 16:52:21.065527+07	2026-05-28 16:53:23.776639+07
53	10	30	UV Print	60000.00	2026-05-28 16:52:34.048312+07	2026-05-28 16:53:33.695796+07
54	10	31	1 x 1 m	10000.00	2026-05-28 16:55:38.541528+07	2026-05-28 16:55:38.541528+07
55	10	31	2 x 1 m	20000.00	2026-05-28 16:56:27.301708+07	2026-05-28 16:56:27.301708+07
56	10	31	2 x 2 m	20000.00	2026-05-28 16:56:47.751078+07	2026-05-28 16:56:47.751078+07
57	10	31	3 x 1 m	25000.00	2026-05-28 16:57:19.511377+07	2026-05-28 16:57:19.511377+07
58	10	31	3 x 2 m	28000.00	2026-05-28 16:57:33.059788+07	2026-05-28 16:57:33.059788+07
59	10	32	4 Titik mata ayam	12000.00	2026-05-28 17:00:32.192204+07	2026-05-28 17:00:32.192204+07
60	10	32	Lipat pinggir	10000.00	2026-05-28 17:00:45.410443+07	2026-05-28 17:00:45.410443+07
61	10	32	Tali	10000.00	2026-05-28 17:00:55.435575+07	2026-05-28 17:00:55.435575+07
62	10	32	Ring	10000.00	2026-05-28 17:01:02.308789+07	2026-05-28 17:01:02.308789+07
63	10	32	Laminasi	25000.00	2026-05-28 17:01:18.537607+07	2026-05-28 17:01:18.537607+07
\.


--
-- Data for Name: service_spesifications; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.service_spesifications (id, service_id, name, input_type, options, is_required, is_active, created_at, updated_at) FROM stdin;
1	1	Jenis Kertas	select	["Sidu", "Natural"]	t	t	2026-01-10 13:57:33.389325+07	2026-01-10 13:57:33.389325+07
2	1	Warna Cetak	select	["Hitam-putih", "Warna"]	t	t	2026-01-10 13:58:47.032746+07	2026-01-10 13:58:47.032746+07
3	1	Ukuran Kertas	select	["A4", "F4"]	t	t	2026-01-10 13:59:45.458356+07	2026-01-10 13:59:45.458356+07
4	1	Jumlah Halaman	number	\N	t	t	2026-01-10 14:00:18.746473+07	2026-01-10 14:00:18.746473+07
6	2	Warna Cetak	select	["Warna", "Hitam-putih"]	t	t	2026-01-12 11:13:15.065866+07	2026-01-12 11:13:15.065866+07
7	2	Ukuran Kertas	select	["A4", "F4"]	t	t	2026-01-12 11:13:56.717762+07	2026-01-12 11:13:56.717762+07
8	2	Jenis Kertas	select	["Sidu", "Natural"]	t	t	2026-01-12 11:18:53.699091+07	2026-01-12 11:18:53.699091+07
9	3	Jenis Kertas	select	["Sidu", "Natural"]	t	t	2026-01-12 11:25:19.409808+07	2026-01-12 11:25:19.409808+07
10	3	Ukuran Kertas	select	["A4", "F4"]	t	t	2026-01-12 11:25:38.967803+07	2026-01-12 11:25:38.967803+07
11	3	Warna Cetak	select	["Warna", "Hitam-putih"]	t	t	2026-01-12 11:25:56.412068+07	2026-01-12 11:25:56.412068+07
12	3	Cetak Depan Belakang	boolean	\N	f	t	2026-01-12 11:28:18.669301+07	2026-01-12 11:28:18.669301+07
13	2	Cetak Depan Belakang	boolean	\N	f	t	2026-01-12 11:28:25.136819+07	2026-01-12 11:28:25.136819+07
14	4	Warna Cetak	select	["Warna", "Hitam-putih"]	t	t	2026-01-12 12:11:51.222027+07	2026-01-12 12:11:51.222027+07
15	4	Ukuran Kertas	select	["A4", "F4"]	t	t	2026-01-12 12:12:13.093178+07	2026-01-12 12:12:13.093178+07
16	4	Jenis Kertas	select	["Sidu", "Natural"]	t	t	2026-01-12 12:12:36.194225+07	2026-01-12 12:12:36.194225+07
17	4	Jumlah Halaman	number	\N	t	t	2026-01-12 12:13:03.965405+07	2026-01-12 12:13:03.965405+07
19	5	Warna background	text	\N	t	t	2026-02-19 14:38:14.596558+07	2026-02-19 14:38:14.596558+07
20	6	Ukuran kertas	select	["A4", "F4"]	t	t	2026-03-31 17:27:01.098796+07	2026-03-31 17:27:01.098796+07
21	6	Bahan	select	["Kertas"]	t	t	2026-03-31 17:28:37.141624+07	2026-03-31 17:28:37.141624+07
22	6	Bentuk	text	\N	t	t	2026-03-31 17:29:32.865474+07	2026-03-31 17:29:32.865474+07
23	6	Finishing	boolean	\N	t	t	2026-03-31 17:30:14.359804+07	2026-03-31 17:30:14.359804+07
24	7	Bahan	select	["Undangan Fadhil RS"]	t	t	2026-05-01 11:07:02.995094+07	2026-05-01 11:07:02.995094+07
25	7	Perlangkapan	select	["Tanpa Perlengkapan", "Label Nama", "Plastik Undangan", "Label Nama dan Plastik"]	t	t	2026-05-01 11:10:37.143899+07	2026-05-01 11:10:37.143899+07
26	8	Bahan	select	["Undangan Maliq 98"]	t	t	2026-05-01 15:39:03.478112+07	2026-05-01 15:39:03.478112+07
27	8	Perlengkapan	select	["Tanpa Perlengkapan", "Label Nama", "Plastik Undangan", "Label Nama dan Plastik"]	t	t	2026-05-05 12:43:47.106493+07	2026-05-05 12:43:47.106493+07
28	9	Bahan	select	["mini banner 25 x 40 - Art Carton", "mini banner 25 x 40 - Art Carton Glossy", "mini banner 25 x 40 - Art Carton Doft"]	t	t	2026-05-21 14:30:39.096871+07	2026-05-21 14:30:39.096871+07
29	5	Paket	select	["sepaket 2x3 isi 6", "sepaket 3x4 isi 4", "sepaket 4x6 isi 4"]	t	t	2026-05-23 12:23:37.766254+07	2026-05-23 12:23:37.766254+07
18	5	Ukuran foto	select	["2x3", "3x4", "4x6"]	f	f	2026-02-19 14:33:32.412872+07	2026-02-19 14:33:32.412872+07
30	10	Bahan	select	["Flexi China", "Flexi Korea", "Flexi Jerman", "UV Print"]	t	t	2026-05-28 16:51:38.125488+07	2026-05-28 16:51:38.125488+07
31	10	Ukuran	select	["1 x 1 m", "2 x 1 m", "2 x 2 m", "3 x 1 m", "3 x 2 m"]	t	t	2026-05-28 16:55:13.245725+07	2026-05-28 16:55:13.245725+07
32	10	Finishing	select	["4 Titik mata ayam", "Lipat pinggir", "Tali", "Ring", "Laminasi"]	f	t	2026-05-28 17:00:14.669185+07	2026-05-28 17:00:14.669185+07
\.


--
-- Data for Name: services; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.services (id, category_id, name, slug, description, base_price, is_active, created_at, updated_at, duration_per_unit) FROM stdin;
4	1	Curiculum Vitae	curiculum-vitae	Layanan print CV ( Curiculum Vitae )	1000.00	t	2026-01-12 11:38:56.188216+07	2026-01-12 11:38:56.188216+07	0
1	1	Print Dokumen	print-dokumen	Layanan print dokumen	1000.00	t	2026-01-09 20:41:22.366531+07	2026-01-09 20:41:22.366531+07	5
2	1	Ijazah	ijazah	Layanan print ijazah	1000.00	t	2026-01-12 11:10:59.521501+07	2026-01-12 11:10:59.521501+07	5
3	1	Sertifikat	sertifikat	Layanan print sertifikat	1000.00	t	2026-01-12 11:17:17.116246+07	2026-01-12 11:17:17.116246+07	5
7	6	Undangan Fadhil RS	undangan-fadhil-rs	Undangan Template\r\n\r\n1. Kamu cukup kirim Nama - Hari & Tanggal saja\r\n2. Pilihan sesuai Gambar produk yg kamu pilih\r\n3. Untuk Tampilan tertutup kamu bisa lihat di tab Foto di kanan atas	1500.00	t	2026-05-01 10:49:33.077613+07	2026-05-01 10:49:33.077613+07	10
8	6	Undangan Malik 98	undangan-malik-98	Ini undangan hasil desain mas malik, dengan ketentuan minimal order 100, dan cukup kasih gambar berdua, sama tanggal dan nama	1000.00	t	2026-05-01 14:38:42.476519+07	2026-05-01 14:38:42.476519+07	15
6	5	Sticker Kertas	sticker-kertas	Sticker berbahan kertas dengan harga lebih ekonomis. Cocok untuk label produk, kemasan, dan kebutuhan indoor.	2000.00	f	2026-03-31 17:25:05.603632+07	2026-03-31 17:25:05.603632+07	0
5	3	Cetak foto	cetak-foto	Menyediakan layanan cetak foto berkualitas tinggi dengan warna tajam dan detail jelas menggunakan kertas foto premium. Tersedia berbagai pilihan ukuran, cocok untuk kebutuhan pribadi, dokumentasi, maupun bisnis, dengan proses cepat dan harga terjangkau.	1000.00	t	2026-02-17 19:01:02.904329+07	2026-02-17 19:01:02.904329+07	5
9	2	Mini Banner	mini-banner	Mini Banner dengan stand X mini\r\nUkuran 25 x 40 cm\r\nBahan Art Carton 260\r\nCetak Full Color	20000.00	t	2026-05-21 14:26:39.475804+07	2026-05-21 14:26:39.475804+07	60
10	2	MMT | Spanduk	mmt-spanduk	Spesifikasi Produk\r\nBahan Flexy Fronlite dengan pilihan ketebalan yg bervariasi\r\n\r\nPrint Outdoor\r\n- Cocok untuk di luar ruangan\r\n- Dicetak dengan resolusi standart\r\n\r\nPrint Indoor \r\n- Cocok untuk di dalam ruangan\r\n- Dicetak menggunakan resolusi tinggi\r\n- Hasil lebih detail jika di bandingkan dengan Print Outdoor\r\n\r\nPrint UV :\r\n- Cocok untuk di dalam ruangan maupun di luar ruangan\r\n- Dicetak menggunakan resolusi sangat tinggi\r\n- Hasil lebih detail jika di bandingkan dengan Print Indoor maupun print outdoor\r\n- Tahan terhadap sinar matahari dan juga tahan terhadap goresan\r\n\r\nCatatan Penting tentang Warna\r\nWarna hasil cetak mungkin sedikit berbeda dengan tampilan di layar Anda. Hal ini normal karena layar menggunakan sistem warna RGB (Red, Green, Blue) berbasis cahaya, sedangkan printer menggunakan tinta CMYK (Cyan, Magenta, Yellow, Black).	30000.00	t	2026-05-28 16:48:53.498467+07	2026-05-28 16:48:53.498467+07	1440
\.


--
-- Name: service_categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.service_categories_id_seq', 6, true);


--
-- Name: service_media_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.service_media_id_seq', 88, true);


--
-- Name: service_specification_values_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.service_specification_values_id_seq', 63, true);


--
-- Name: service_specifications_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.service_specifications_id_seq', 32, true);


--
-- Name: services_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.services_id_seq', 10, true);


--
-- Name: service_categories service_categories_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_categories
    ADD CONSTRAINT service_categories_name_key UNIQUE (name);


--
-- Name: service_categories service_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_categories
    ADD CONSTRAINT service_categories_pkey PRIMARY KEY (id);


--
-- Name: service_categories service_categories_slug_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_categories
    ADD CONSTRAINT service_categories_slug_key UNIQUE (slug);


--
-- Name: service_media service_media_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_media
    ADD CONSTRAINT service_media_pkey PRIMARY KEY (id);


--
-- Name: service_spesification_values service_specification_values_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_spesification_values
    ADD CONSTRAINT service_specification_values_pkey PRIMARY KEY (id);


--
-- Name: service_spesifications service_specifications_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_spesifications
    ADD CONSTRAINT service_specifications_pkey PRIMARY KEY (id);


--
-- Name: services services_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.services
    ADD CONSTRAINT services_pkey PRIMARY KEY (id);


--
-- Name: services services_slug_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.services
    ADD CONSTRAINT services_slug_key UNIQUE (slug);


--
-- Name: service_media service_media_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_media
    ADD CONSTRAINT service_media_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.services(id) ON DELETE CASCADE;


--
-- Name: service_spesification_values service_specification_values_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_spesification_values
    ADD CONSTRAINT service_specification_values_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.services(id) ON DELETE CASCADE;


--
-- Name: service_spesification_values service_specification_values_specification_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_spesification_values
    ADD CONSTRAINT service_specification_values_specification_id_fkey FOREIGN KEY (spesification_id) REFERENCES public.service_spesifications(id) ON DELETE CASCADE;


--
-- Name: service_spesifications service_specifications_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_spesifications
    ADD CONSTRAINT service_specifications_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.services(id) ON DELETE CASCADE;


--
-- Name: services services_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.services
    ADD CONSTRAINT services_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.service_categories(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict EiM8dBeMW4rh5hRrcCbqBJQtAXpO8nHq8udUQQR2nURXwo9BtxzI3t1nv3cvVny

