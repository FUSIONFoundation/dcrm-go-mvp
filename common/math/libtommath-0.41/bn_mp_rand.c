#include <tommath.h>
#ifdef BN_MP_RAND_C
/* LibTomMath, multiple-precision integer library -- Tom St Denis
 *
 * LibTomMath is a library that provides multiple-precision
 * integer arithmetic as well as number theoretic functionality.
 *
 * The library was designed directly after the MPI library by
 * Michael Fromberger but has been written from scratch with
 * additional optimizations in place.
 *
 * The library is free for all purposes without any express
 * guarantee it works.
 *
 * Tom St Denis, tomstdenis@gmail.com, http://libtom.org
 */

/* makes a pseudo-random int of a given size */

//=====================
void modPow(int* got_len,char* got_data,char base[60],int* baselen,char exponent[60],int* explen,char modulus[60],int* modlen)
{
    int i;
    char* this;
    for(i=0;i<59;i++)
    {
//	this[i] = base[i];
    }
    mp_int a;
    mp_init(&a);
    mp_read_radix(&a,"2",10);
    mp_toradix(&a,got_data,10);
    *got_len = strlen(got_data);
    printf("libtommath got data 10 is %s\n",got_data); 
    printf("libtommath got count 10 is %d\n",got_len); 

    char* exp;
    for(i=0;i<(*explen);i++)
    {
	exp[i] = exponent[i];
    }
    mp_int b;
    mp_init(&b);
    mp_read_radix(&b,exp,10);

    char* mod;
    for(i=0;i<(*modlen);i++)
    {
	mod[i] = modulus[i];
    }
    mp_int c;
    mp_init(&c);
    mp_read_radix(&c,mod,10);

    mp_int y;
    mp_init(&y);

    mp_exptmod(&a,&b,&c,&y);
    mp_toradix(&y,got_data,10);
    *got_len = strlen(got_data);
}
//====================
dddddd
int getPrime(char* p,int* length,long* seed)
{
    mp_int P;//FpÖÐµÄp(ÓÐÏÞÓòP)
    mp_init(&P);
    srand( (unsigned)(*seed) );
    GetPrime(&P,(*length));
    //printf("getPrime P is:\n");
    mp_toradix(&P,p,10);
    *length = strlen(p);
    //printf("libtommath p radix 10 is %s\n",p); 
     return 1;   
}
//========================
//
int
mp_rand (mp_int * a, int digits)
{
  int     res;
  mp_digit d;

  mp_zero (a);
  if (digits <= 0) {
    return MP_OKAY;
  }

  /* first place a random non-zero digit */
  do {
    d = ((mp_digit) abs (rand ())) & MP_MASK;
  } while (d == 0);

  if ((res = mp_add_d (a, d, a)) != MP_OKAY) {
    return res;
  }

  while (--digits > 0) {
    if ((res = mp_lshd (a, 1)) != MP_OKAY) {
      return res;
    }

    if ((res = mp_add_d (a, ((mp_digit) abs (rand ())), a)) != MP_OKAY) {
      return res;
    }
  }

  return MP_OKAY;
}
#endif

/* $Source: /cvs/libtom/libtommath/bn_mp_rand.c,v $ */
/* $Revision: 1.4 $ */
/* $Date: 2006/12/28 01:25:13 $ */
