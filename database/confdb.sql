--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.7
-- Dumped by pg_dump version 9.6.7

-- Started on 2018-02-14 10:25:26 EET

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

DROP DATABASE confdb;
--
-- TOC entry 2169 (class 1262 OID 315920)
-- Name: confdb; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE confdb WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8';


ALTER DATABASE confdb OWNER TO postgres;

\connect confdb

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 1 (class 3079 OID 12427)
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- TOC entry 2172 (class 0 OID 0)
-- Dependencies: 1
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- TOC entry 185 (class 1259 OID 315921)
-- Name: mongodb; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE mongodb (
    domain character(1),
    mongodb boolean,
    host character(1),
    port character(1)
);


ALTER TABLE mongodb OWNER TO postgres;

--
-- TOC entry 187 (class 1259 OID 315927)
-- Name: tempconfig; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE tempconfig (
    rest_api_root character(1),
    host character(1),
    port character(1),
    remoting character(1),
    legasy_explorer boolean
);


ALTER TABLE tempconfig OWNER TO postgres;

--
-- TOC entry 186 (class 1259 OID 315924)
-- Name: tsconfig; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE tsconfig (
    module character(1),
    target character(1),
    source_map boolean,
    exclude integer
);


ALTER TABLE tsconfig OWNER TO postgres;

--
-- TOC entry 2162 (class 0 OID 315921)
-- Dependencies: 185
-- Data for Name: mongodb; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- TOC entry 2164 (class 0 OID 315927)
-- Dependencies: 187
-- Data for Name: tempconfig; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- TOC entry 2163 (class 0 OID 315924)
-- Dependencies: 186
-- Data for Name: tsconfig; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- TOC entry 2171 (class 0 OID 0)
-- Dependencies: 7
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

GRANT ALL ON SCHEMA public TO PUBLIC;


-- Completed on 2018-02-14 10:25:26 EET

--
-- PostgreSQL database dump complete
--

