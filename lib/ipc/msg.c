/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
#include "msg.h"

int msgGet(key_t key, int msgflg)
{
	return msgget(key, msgflg);
}

int msgSnd(int msqid, const void *msgp, size_t msgsz, int msgtyp, int msgflg)
{
	struct msgbuf m;

	if (msgsz > MSGSZ) {
		return -1;
	}

	memcpy(m.mtext, msgp, msgsz);
	m.mtype = msgtyp;

	return msgsnd(msqid, &m, msgsz, msgflg);
}

int msgRcv(int msqid, int msgtyp, int msgflg, char **pp)
{
	struct msgbuf m;
	int len;

	len = msgrcv(msqid, &m, MSGSZ, 0, 0);
	if (len < 0) {
		return len;
	}

	*pp = (char *)malloc(len);
	memcpy(*pp, m.mtext, len);
	return len;

	//return len;
}
