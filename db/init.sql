CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;
COMMENT ON EXTENSION citext IS 'data type for case-insensitive character strings';
CREATE FUNCTION public.insert_post() RETURNS trigger
    LANGUAGE plpgsql
AS $$
DECLARE
    parent_path         bigint[];
    first_parent_thread integer;
BEGIN
    IF NEW.parent IS NULL THEN
        NEW.path := array_append(NEW.path, NEW.id);
    ELSE
        SELECT path FROM posts WHERE id = NEW.parent INTO parent_path;
        SELECT thread FROM posts WHERE id = parent_path[1] INTO first_parent_thread;
        IF NOT FOUND OR first_parent_thread <> NEW.thread THEN
            RAISE EXCEPTION 'Parent post was created in another thread' USING ERRCODE = '666';
        END IF;
        NEW.path := NEW.path || parent_path || NEW.id;
    END IF;
    UPDATE forums SET posts = posts + 1 WHERE forums.slug = NEW.forum;
    RETURN NEW;
END
$$;
ALTER FUNCTION public.insert_post() OWNER TO test;
CREATE FUNCTION public.count_forum_posts() RETURNS trigger
    LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE forums
    SET posts = posts + 1
    WHERE slug = NEW.forum;
    RETURN NULL;
END;
$$;
ALTER FUNCTION public.count_forum_posts() OWNER TO test;
CREATE FUNCTION public.count_forum_threads() RETURNS trigger
    LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE forums
    SET threads = threads + 1
    WHERE slug = NEW.forum;
    RETURN NULL;
END;
$$;
ALTER FUNCTION public.count_forum_threads() OWNER TO test;
CREATE FUNCTION public.path() RETURNS trigger
    LANGUAGE plpgsql
AS $$
DECLARE
    parent_path      integer[];
    parent_thread_id integer;
BEGIN
    IF NEW.parent is NULL THEN
        NEW.path := NEW.path || NEW.id;
    ELSE
        SELECT path, thread
        FROM posts
        WHERE id = NEW.parent
        INTO parent_path, parent_thread_id;
        IF parent_thread_id != NEW.thread THEN
            raise exception 'Path error';
        end if;
        NEW.path := NEW.path || parent_path || NEW.id;
    END IF;
    RETURN NEW;
END;
$$;
ALTER FUNCTION public.path() OWNER TO test;
CREATE FUNCTION public.edit_post() RETURNS trigger
    LANGUAGE plpgsql
AS $$
BEGIN
    IF NEW.message = OLD.message
    THEN RETURN NULL;
    END IF;
    UPDATE posts SET isedited = TRUE
    WHERE id = NEW.id;
    RETURN NULL;
END;
$$;
ALTER FUNCTION public.edit_post() OWNER TO test;
CREATE FUNCTION public.vote_post() RETURNS trigger
    LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE threads
    SET votes = votes + NEW.voice
    WHERE id = NEW.thread;
    RETURN NULL;
END;
$$;
ALTER FUNCTION public.vote_post() OWNER TO test;
CREATE FUNCTION public.vote_update() RETURNS trigger
    LANGUAGE plpgsql
AS $$
BEGIN
    IF OLD.voice = NEW.voice
    THEN
        RETURN NULL;
    END IF;
    UPDATE threads
    SET
        votes = votes + CASE WHEN NEW.voice = -1
                                 THEN -2
                             ELSE 2 END
    WHERE id = NEW.thread;
    RETURN NULL;
END;
$$;
ALTER FUNCTION public.vote_update() OWNER TO test;
SET default_tablespace = '';
SET default_table_access_method = heap;
CREATE UNLOGGED TABLE public.forum_users (
                                             user_nickname public.citext NOT NULL COLLATE pg_catalog."C",
                                             forum_slug public.citext NOT NULL COLLATE pg_catalog."C"
);
ALTER TABLE public.forum_users OWNER TO test;
CREATE UNLOGGED TABLE public.forums (
                                        id integer NOT NULL,
                                        title text NOT NULL,
                                        "user" public.citext NOT NULL COLLATE pg_catalog."C",
                                        slug public.citext NOT NULL COLLATE pg_catalog."C",
                                        posts bigint DEFAULT 0 NOT NULL,
                                        threads integer DEFAULT 0 NOT NULL
);
ALTER TABLE public.forums OWNER TO test;
CREATE SEQUENCE public.forums_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER TABLE public.forums_id_seq OWNER TO test;
ALTER SEQUENCE public.forums_id_seq OWNED BY public.forums.id;
CREATE UNLOGGED TABLE public.posts (
                                       id bigint NOT NULL,
                                       parent bigint DEFAULT 0,
                                       author public.citext NOT NULL COLLATE pg_catalog."C",
                                       message text NOT NULL,
                                       isedited boolean DEFAULT false,
                                       forum public.citext COLLATE pg_catalog."C",
                                       thread integer,
                                       created timestamp with time zone DEFAULT now(),
                                       path bigint[]
);
ALTER TABLE public.posts OWNER TO test;
CREATE SEQUENCE public.posts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER TABLE public.posts_id_seq OWNER TO test;
ALTER SEQUENCE public.posts_id_seq OWNED BY public.posts.id;
CREATE UNLOGGED TABLE public.threads (
                                         id integer NOT NULL,
                                         title public.citext NOT NULL COLLATE pg_catalog."C",
                                         author public.citext NOT NULL COLLATE pg_catalog."C",
                                         forum public.citext COLLATE pg_catalog."C",
                                         message text NOT NULL,
                                         votes integer DEFAULT 0,
                                         slug public.citext COLLATE pg_catalog."C",
                                         created timestamp with time zone DEFAULT now()
);
ALTER TABLE public.threads OWNER TO test;
CREATE SEQUENCE public.threads_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER TABLE public.threads_id_seq OWNER TO test;
ALTER SEQUENCE public.threads_id_seq OWNED BY public.threads.id;
CREATE UNLOGGED TABLE public.users (
                                       id integer NOT NULL,
                                       nickname public.citext COLLATE pg_catalog."C",
                                       fullname public.citext NOT NULL COLLATE pg_catalog."C",
                                       about text,
                                       email public.citext NOT NULL COLLATE pg_catalog."C"
);
ALTER TABLE public.users OWNER TO test;
CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER TABLE public.users_id_seq OWNER TO test;
ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;
CREATE UNLOGGED TABLE public.votes (
                                       nickname public.citext NOT NULL COLLATE pg_catalog."C",
                                       voice integer NOT NULL,
                                       thread integer NOT NULL
);
ALTER TABLE public.votes OWNER TO test;
ALTER TABLE ONLY public.forums ALTER COLUMN id SET DEFAULT nextval('public.forums_id_seq'::regclass);
ALTER TABLE ONLY public.posts ALTER COLUMN id SET DEFAULT nextval('public.posts_id_seq'::regclass);
ALTER TABLE ONLY public.threads ALTER COLUMN id SET DEFAULT nextval('public.threads_id_seq'::regclass);
ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);
ALTER TABLE ONLY public.forums
    ADD CONSTRAINT forum_pk PRIMARY KEY (id);
ALTER TABLE ONLY public.posts
    ADD CONSTRAINT post_pk PRIMARY KEY (id);
ALTER TABLE ONLY public.threads
    ADD CONSTRAINT thread_pk PRIMARY KEY (id);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT user_pk PRIMARY KEY (id);
ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_user_thread_unique UNIQUE (thread, nickname);


CREATE UNIQUE INDEX uidx_forum_users_user_id_forum_id ON public.forum_users USING btree (user_nickname, forum_slug);
CREATE UNIQUE INDEX uidx_forums_slug ON public.forums USING btree (slug);
CREATE UNIQUE INDEX uidx_forums_user ON public.forums USING btree ("user");
CREATE UNIQUE INDEX uidx_users_nickname ON public.users USING btree (nickname);
CREATE UNIQUE INDEX uidx_threads_slug ON public.threads USING btree (slug);
CREATE UNIQUE INDEX uidx_user_id ON public.users USING btree (id);
CREATE UNIQUE INDEX uidx_users_email ON public.users USING btree (email);

CREATE INDEX idx_post_threadid_created_id ON public.posts USING btree (thread, created, id, parent, path);
-- CREATE INDEX idx_post_threadid_path ON public.posts USING btree (thread, path);
-- CREATE INDEX idx_posts_id ON public.posts USING hash (id);

-- CREATE INDEX idx_threads_slug_hash ON public.threads USING hash (slug);
CREATE INDEX idx_threads_forum_created ON public.threads USING btree (forum, created);

-- CREATE INDEX idx_forums_slug_hash ON public.forums USING hash (slug);
-- CREATE INDEX idx_forums_users_foreign ON public.forums USING hash ("user");


CREATE TRIGGER before_insert_post BEFORE INSERT ON public.posts FOR EACH ROW EXECUTE FUNCTION public.insert_post();
CREATE TRIGGER count_forum_threads AFTER INSERT ON public.threads FOR EACH ROW EXECUTE FUNCTION public.count_forum_threads();
CREATE TRIGGER edit_post AFTER UPDATE ON public.posts FOR EACH ROW EXECUTE FUNCTION public.edit_post();
CREATE TRIGGER vote_post AFTER INSERT ON public.votes FOR EACH ROW EXECUTE FUNCTION public.vote_post();
CREATE TRIGGER vote_update AFTER UPDATE ON public.votes FOR EACH ROW EXECUTE FUNCTION public.vote_update();
ALTER TABLE ONLY public.forum_users
    ADD CONSTRAINT forum_users_forum_slug_fk FOREIGN KEY (forum_slug) REFERENCES public.forums(slug);
ALTER TABLE ONLY public.forum_users
    ADD CONSTRAINT forum_users_user_nickname_fk FOREIGN KEY (user_nickname) REFERENCES public.users(nickname);
ALTER TABLE ONLY public.forums
    ADD CONSTRAINT forums_users_nickname_fk FOREIGN KEY ("user") REFERENCES public.users(nickname);
ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_users_nickname_fk FOREIGN KEY (author) REFERENCES public.users(nickname);
ALTER TABLE ONLY public.threads
    ADD CONSTRAINT threads_users_nickname_fk FOREIGN KEY (author) REFERENCES public.users(nickname);
ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_users_nickname_fk FOREIGN KEY (nickname) REFERENCES public.users(nickname);

