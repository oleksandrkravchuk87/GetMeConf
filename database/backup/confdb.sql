--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.7
-- Dumped by pg_dump version 9.6.7

-- Started on 2018-02-19 13:10:58 EET

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
-- TOC entry 2172 (class 1262 OID 315920)
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
-- TOC entry 2175 (class 0 OID 0)
-- Dependencies: 1
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- TOC entry 185 (class 1259 OID 315921)
-- Name: mongodbs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE mongodbs (
    domain text,
    mongodb boolean,
    host text,
    port text
);


ALTER TABLE mongodbs OWNER TO postgres;

--
-- TOC entry 186 (class 1259 OID 315927)
-- Name: tempconfigs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE tempconfigs (
    rest_api_root text,
    host text,
    port text,
    remoting text,
    legasy_explorer boolean
);


ALTER TABLE tempconfigs OWNER TO postgres;

--
-- TOC entry 187 (class 1259 OID 315994)
-- Name: tsconfigs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE tsconfigs (
    module text,
    target text,
    source_map boolean,
    excluding integer
);


ALTER TABLE tsconfigs OWNER TO postgres;

--
-- TOC entry 2165 (class 0 OID 315921)
-- Dependencies: 185
-- Data for Name: mongodbs; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO mongodbs (domain, mongodb, host, port) VALUES ('mydom', true, 'localhost', '8080');
INSERT INTO mongodbs (domain, mongodb, host, port) VALUES ('testdom', true, '127.0.0.1', '9090');
INSERT INTO mongodbs (domain, mongodb, host, port) VALUES ('remote', true, '227.255.255.1', '8090');
INSERT INTO mongodbs (domain, mongodb, host, port) VALUES ('asia', true, '217.155.155.1', '8081');


--
-- TOC entry 2166 (class 0 OID 315927)
-- Dependencies: 186
-- Data for Name: tempconfigs; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO tempconfigs (rest_api_root, host, port, remoting, legasy_explorer) VALUES ('root', '', '', '', false);
INSERT INTO tempconfigs (rest_api_root, host, port, remoting, legasy_explorer) VALUES ('rest', 'asia', '8080', 'local_uk', true);
INSERT INTO tempconfigs (rest_api_root, host, port, remoting, legasy_explorer) VALUES ('api', 'europa', '9080', 'local', true);
INSERT INTO tempconfigs (rest_api_root, host, port, remoting, legasy_explorer) VALUES ('local_', 'localhost', '8080', 'rem', true);


--
-- TOC entry 2167 (class 0 OID 315994)
-- Dependencies: 187
-- Data for Name: tsconfigs; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO tsconfigs (module, target, source_map, excluding) VALUES ('admin', 'admins', true, 1);
INSERT INTO tsconfigs (module, target, source_map, excluding) VALUES ('user', 'users', true, 1);
INSERT INTO tsconfigs (module, target, source_map, excluding) VALUES ('customer', 'customers', true, 100);
INSERT INTO tsconfigs (module, target, source_map, excluding) VALUES ('vendor', 'vendors', true, 33);


--
-- TOC entry 2174 (class 0 OID 0)
-- Dependencies: 7
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

GRANT ALL ON SCHEMA public TO PUBLIC;


-- Completed on 2018-02-19 13:10:58 EET

--
-- PostgreSQL database dump complete
--

