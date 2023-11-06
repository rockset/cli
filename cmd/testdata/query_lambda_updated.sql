SELECT
    COUNT(e.type) AS type_cnt,
    COUNT(e.kind) AS kind_cnt,
    e.label
FROM
    commons._events e
GROUP BY
    e.label
ORDER BY
    type_cnt DESC,
    kind_cnt DESC,
    e.label
