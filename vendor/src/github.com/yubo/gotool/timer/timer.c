/*
 * yubo@yubo.org
 * 2015-12-04
 */
#include "timer.h"

int64_t nanotime(void){
	struct timespec ts; 
	clock_gettime(CLOCK_MONOTONIC, &ts);
	return ts.tv_sec*1000000000+ts.tv_nsec;
}
