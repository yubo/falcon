/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

#ifndef __CACHE_H__
#define __CACHE_H__

#include <stdlib.h>
#include <stdio.h>
#include "shm.h"

#define CACHE_SIZE (1<<5)
#define C_HASH_SIZE 32
#define C_HOST_SIZE 32
#define C_NAME_SIZE 32
#define C_TAGS_SIZE 32
#define C_TYPE_SIZE 32

struct rrd_data {
	int64_t time;
	double value;
};

struct cache_entry {
	uint32_t flag;
	char hashkey[C_HASH_SIZE+1];
	int64_t idxTs;
	int64_t commitTs;
	int64_t createTs;
	int64_t lastTs;
	char host[C_HOST_SIZE+1];
	char name[C_NAME_SIZE+1];
	char tags[C_TAGS_SIZE+1];
	char typ[C_TYPE_SIZE+1];
	int step;
	int heartbeat;
	char min;
	char max;
	uint32_t dataId;
	uint32_t commitId;
	int64_t  time[CACHE_SIZE];
	double value[CACHE_SIZE];
};

struct cache_block {
	uint32_t magic;
	int block_size;
	int cache_entry_nb;
	int rrd_data_nb;
};

static inline int cache_entry_reset(struct cache_entry *e, int64_t createTs,
		char *host, char *name, char *tags, char *typ,
		int step, int heartbeat, char min, char max){
	if( strlen(host) > C_HOST_SIZE || 
			strlen(name) > C_NAME_SIZE ||
			strlen(tags) > C_TAGS_SIZE ||
			strlen(typ) > C_TYPE_SIZE)
		return 1;
	memset(e, 0, sizeof(*e));
	strcpy(&e->host[0], host);
	strcpy(&e->name[0], name);
	strcpy(&e->tags[0], tags);
	strcpy(&e->typ[0], typ);
	e->createTs = createTs;
	e->step = step;
	e->heartbeat = heartbeat;
	e->min = min;
	e->max = max;
	return 0;
}

static inline int set_hashkey(struct cache_entry *e, char *hashkey){
	return snprintf(&e->hashkey[0], C_HASH_SIZE, "%s", hashkey);
}

static inline int set_host(struct cache_entry *e, char *host){
	return snprintf(&e->host[0], C_HOST_SIZE, "%s", host);
}

static inline int set_name(struct cache_entry *e, char *name){
	return snprintf(&e->name[0], C_NAME_SIZE, "%s", name);
}

static inline int set_tags(struct cache_entry *e, char *tags){
	return snprintf(&e->tags[0], C_TAGS_SIZE, "%s", tags);
}

static inline int set_typ(struct cache_entry *e, char *typ){
	return snprintf(&e->typ[0], C_TYPE_SIZE, "%s", typ);
}


#endif




