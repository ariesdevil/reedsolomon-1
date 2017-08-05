#include "textflag.h"

#define tbl AX
#define datas  BX
#define parity CX
#define num_data SI
#define len DI
#define done R8
#define mask Y2

#define tmp0 R9
#define tmp1 R10
#define tmp2 R11
#define tmp3 R12


// func vectMul(tbl []byte, in Matrix, out []byte)
TEXT ·vectMul(SB), NOSPLIT, $0
	// restore args
	MOVQ  tblAddr+0(FP), tbl
	MOVQ  in+24(FP), datas
	MOVQ  out+48(FP), parity
	MOVQ  size+56(FP), len

	// init mask
	MOVB         $0x0f, DX
    LONG         $0x2069e3c4; WORD $0x00d2   // VPINSRB $0x00, EDX, XMM2, XMM2
   	VPBROADCASTB X2, mask

   	// init done
    MOVQ  $0, done

loop:
	// init tbl
	MOVQ         tbl, tmp0      // restore addr of tbl
	VMOVDQU      (tmp0), X0       // get lowtbl
    VMOVDQU      16(tmp0), X1     // get hightbl
    VINSERTI128  $1, X0, Y0, Y0
   	VINSERTI128  $1, X1, Y1, Y1

   	MOVQ num_data, tmp1
   	SUBQ $2, tmp1
   	MOVQ $0, tmp2
   	MOVQ (datas)(tmp2*1), tmp3
   	// maybe non-tmp get
   	VMOVDQU (tmp3)(done*1), Y3
   	VPSRLQ  $4, Y3, Y4
   	VPAND   mask, Y4, Y4
   	VPAND   mask, Y3, Y3
   	VPSHUFB Y4, Y1, Y4
   	VPSHUFB Y3, Y0, Y3
   	VPSHUFB Y3, Y4, Y5

next_vect:
	ADDQ  $32, tmp0
	VMOVDQU      (tmp0), X0       // get lowtbl
    VMOVDQU      16(tmp0), X1     // get hightbl
    VINSERTI128  $1, X0, Y0, Y0
    VINSERTI128  $1, X1, Y1, Y1

	ADDQ  $24, tmp2
	MOVQ  (datas)(tmp2*1), tmp3
	VMOVDQU (tmp3)(done*1), Y3
	VPSRLQ  $4, Y3, Y4
	VPAND   mask, Y4, Y4
	VPAND   mask, Y3, Y3
	VPSHUFB Y4, Y1, Y4
	VPSHUFB Y3, Y0, Y3
	VPXOR   Y3, Y4, Y6
	VPXOR   Y5, Y6, Y5
	SUBQ    $1, tmp1
	JGE     next_vect

	VMOVNTDQ Y5, (parity)(done*1)
	ADDQ   $32, done
	CMPQ   len, done
	JNE    loop
	RET
