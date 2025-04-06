-- api/videos.go
-- name: GetVideos :many
select 
    sqlc.embed(videos), sqlc.embed(channels)
from videos
left join channels
on videos.owner = channels.id
order by videos.rowid desc;

-- name: GetVideo :one
select
    sqlc.embed(videos), sqlc.embed(channels)
from videos
left join channels
on videos.owner = channels.id
where videos.id=?;

-- api/channels.go
-- name: GetChannels :many
select * from channels;

-- name: GetChannel :one
select * from channels
where channels.id=?;

-- name: GetChannelVideos :many
select
    sqlc.embed(v), sqlc.embed(c)
from videos as v
left join channels as c
on v.owner = c.id
where v.owner=?
order by v.rowid desc;

-- api/playlists.go
-- name: GetPlaylists :many
select 
    sqlc.embed(p), sqlc.embed(c), sqlc.embed(v)
from playlists as p
left join channels as c
on p.owner=c.id
left join playlist_video as pv
on p.id=pv.playlistId
left join videos as v
on pv.videoId=v.id	
where c.thumbnail not null
and v.id not null
order by p.rowid desc, pv.sortIndex asc, v.timestamp asc;

-- name: GetPlaylistVideos :many
select 
    sqlc.embed(p), sqlc.embed(c), sqlc.embed(v), sqlc.embed(vc), pv.sortIndex
from playlists as p
left join channels as c
on p.owner=c.id
left join playlist_video as pv
on p.id=pv.playlistId
left join videos as v
on pv.videoId=v.id	
left join channels as vc
on v.owner=vc.id
where p.id=?
order by pv.sortIndex asc, v.timestamp asc;

-- name: UpdatePlaylistVideoIndex :execrows
update playlist_video 
set sortIndex=? 
where playlistId=? and videoId=?