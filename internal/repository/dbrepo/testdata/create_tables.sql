


CREATE TABLE public.users (
  id integer NOT NULL,
  user_name character varying(255),
  email character varying(255),
  password character varying(255),
  avatar_img character varying(255),
  profile text,
  created_at timestamp without time zone,
  updated_at timestamp without time zone
);



ALTER TABLE public.users ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);



CREATE TABLE public.articles (
  id integer NOT NULL,
  content text,
  title character varying(512),
  main_img character varying(512),
  work character varying(255),
  medium integer,
  comment_ok boolean,
  created_at timestamp without time zone,
  updated_at timestamp without time zone,
  user_id integer 
);



ALTER TABLE public.articles ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.articles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);




CREATE TABLE public.comments (
  id integer NOT NULL,
  comment text,
  created_at timestamp without time zone,
  updated_at timestamp without time zone,
  user_id integer,
  article_id integer 
);



ALTER TABLE public.comments ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.comments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);



CREATE TABLE public.article_imgs (
  id integer NOT NULL,
  article_img character varying(255),
  article_id integer
);



ALTER TABLE public.article_imgs ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.article_imgs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);



CREATE TABLE public.article_goods (
  id integer NOT NULL,
  user_id integer,
  article_id integer
);


ALTER TABLE public.article_goods ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.article_goods_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);



CREATE TABLE public.comment_goods (
  id integer NOT NULL,
  user_id integer,
  comment_id integer
);



ALTER TABLE public.comment_goods ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.comment_goods_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT articles_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.article_imgs
    ADD CONSTRAINT article_imgs_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.article_goods
    ADD CONSTRAINT article_goods_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.comment_goods
    ADD CONSTRAINT comment_goods_pkey PRIMARY KEY (id);


ALTER TABLE ONLY public.articles
    ADD CONSTRAINT articles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_article_id_fkey FOREIGN KEY (article_id) REFERENCES public.articles(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE ONLY public.article_imgs
    ADD CONSTRAINT article_imgs_article_id_fkey FOREIGN KEY (article_id) REFERENCES public.articles(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE ONLY public.article_goods
    ADD CONSTRAINT article_goods_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.article_goods
    ADD CONSTRAINT article_goods_article_id_fkey FOREIGN KEY (article_id) REFERENCES public.articles(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE ONLY public.comment_goods
    ADD CONSTRAINT comment_goods_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.comment_goods
    ADD CONSTRAINT comment_goods_comment_id_fkey FOREIGN KEY (comment_id) REFERENCES public.comments(id) ON UPDATE CASCADE ON DELETE CASCADE;




CREATE TABLE public.genres (
  id integer NOT NULL,
  genre_name character varying(255)
);




ALTER TABLE public.genres ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.genres_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


CREATE TABLE public.genre_articles (
  id integer NOT NULL,
  article_id integer,
  genre_id integer
);



ALTER TABLE public.genre_articles ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.genre_articles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);







ALTER TABLE ONLY public.genres
    ADD CONSTRAINT genres_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.genre_articles
    ADD CONSTRAINT genre_articles_pkey PRIMARY KEY (id);



ALTER TABLE ONLY public.genre_articles
    ADD CONSTRAINT genre_articles_article_id_fkey FOREIGN KEY (article_id) REFERENCES public.articles(id);
ALTER TABLE ONLY public.genre_articles
    ADD CONSTRAINT genre_articles_genre_id_fkey FOREIGN KEY (genre_id) REFERENCES public.genres(id);