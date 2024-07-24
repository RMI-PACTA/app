BEGIN;

ALTER TYPE file_type ADD VALUE 'js.map';
ALTER TYPE file_type ADD VALUE 'woff';
ALTER TYPE file_type ADD VALUE 'woff2';
ALTER TYPE file_type ADD VALUE 'eot';
ALTER TYPE file_type ADD VALUE 'svg';
ALTER TYPE file_type ADD VALUE 'png';
ALTER TYPE file_type ADD VALUE 'jpg';
ALTER TYPE file_type ADD VALUE 'pdf';

COMMIT;
