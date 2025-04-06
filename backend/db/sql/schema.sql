create table if not exists channels (
	id text primary key, 
	title text not null, 
	description text not null, 
	thumbnail text not null
);

create table if not exists videos (
	id text primary key, 
	title text not null, 
	description text not null, 
	timestamp timestamp not null, 
	duration text not null, 
	owner text not null, 
	thumbnail text not null
);

create table if not exists playlists (
	id text primary key unique, 
	title text not null, 
	description text not null, 
	timestamp timestamp not null, 
	owner text not null
);

create table if not exists playlist_video (
	playlistId text not null, 
	videoId text not null, 
	sortIndex integer not null, 
	unique (playlistId, videoId), 
	primary key (playlistId, videoId)
);