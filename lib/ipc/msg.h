/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

#ifndef __IPC_H__
#define __IPC_H__

#include <stdio.h>
#include <stdint.h>
#include <string.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/ipc.h>
#include <sys/msg.h>

#define MSGSZ 256

struct msgbuf {
	long mtype;
	char mtext[MSGSZ];
};

int msgGet(key_t key, int msgflg);
int msgSnd(int msqid, const void *msgp, size_t msgsz,int msgtyp, int msgflg);
int msgRcv(int msqid, int msgtyp, int msgflg, char **pp);
#endif
