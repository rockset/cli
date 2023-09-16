SELECT
    _input._meta,
    NULLIF(_input.belongs_to_collection, '') AS belongs_to_collection,
    _input.genres,
    NULLIF(_input.homepage, '') AS homepage,
    TRY_CAST(_input.id AS int) AS id,
    _input.original_language,
    _input.overview,
    _input.popularity,
    _input.production_companies,
    CAST(_input.release_date AS date) AS release_date,
    _input.revenue,
    _input.runtime,
    _input.title,
    _input.vote_average,
    _input.vote_count
FROM
    _input;
