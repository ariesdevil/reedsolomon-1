#include "textflag.h"

#define data BX
#define parity CX
#define len DI
#define pos SI
#define vect AX
#define tbl  R12

#define mask Y2
#define lowTbl Y0
#define highTbl Y1

#define tmp R8
#define tmp1 R9
#define tmp2 R10
#define tmp3 R11

// func in0mul(tbl []byte, data []byte, parity Matrix)
TEXT ·in0mul(SB), NOSPLIT, $0
	MOVQ t+0(FP), tbl
	MOVQ b+24(FP), data
	MOVQ c+32(FP), len
	MOVQ d+48(FP), parity
	MOVQ e+56(FP), vect
	MOVB         $0x0f, DX
    LONG         $0x2069e3c4; WORD $0x00d2
    VPBROADCASTB X2, Y2
    MOVQ $0, pos

loop32b:
	VMOVDQU (data)(pos*1), Y3
	MOVQ    $0, tmp
	VMOVDQU (tbl)(tmp*1), X0
    VMOVDQU 16(tbl)(tmp*1), X1
    VINSERTI128  $1, X0, Y0, Y0
    VINSERTI128  $1, X1, Y1, Y1
    VPSRLQ  $4, Y3, Y4
    VPAND   Y2, Y4, Y4
    VPAND   Y2, Y3, Y5
    VPSHUFB Y4, Y1, Y3
    VPSHUFB Y5, Y0, Y6
    VPXOR   Y3, Y6, Y3
    MOVQ   $0, tmp1
    MOVQ   (parity)(tmp1*1), tmp2
    VMOVDQU Y3, (tmp2)(pos*1)
    MOVQ   vect, tmp3
    SUBQ   $2, tmp3
    JL ret
next_vect:
	ADDQ   $32, tmp
	VMOVDQU (tbl)(tmp*1), X0
    VMOVDQU 16(tbl)(tmp*1), X1
    VINSERTI128  $1, X0, Y0, Y0
    VINSERTI128  $1, X1, Y1, Y1
	VPSHUFB  Y4, Y1, Y3
	VPSHUFB  Y5, Y0, Y6
	VPXOR Y3, Y6, Y3
	ADDQ   $24, tmp1
	MOVQ   (parity)(tmp1*1), tmp2
	VMOVDQU Y3, (tmp2)(pos*1)
	SUBQ   $1, tmp3
	JGE   next_vect

	ADDQ  $32, pos
	CMPQ  len, pos
	JNZ   loop32b
	RET

ret:
	RET

// func in0mulxor(tbl []byte, data []byte, parity Matrix)
TEXT ·in0mulxor(SB), NOSPLIT, $0
	MOVQ t+0(FP), tbl
	MOVQ b+24(FP), data
	MOVQ c+32(FP), len
	MOVQ d+48(FP), parity
	MOVQ e+56(FP), vect
	MOVB         $0x0f, DX
    LONG         $0x2069e3c4; WORD $0x00d2
    VPBROADCASTB X2, Y2
    MOVQ $0, pos

loop32b:
	VMOVDQU (data)(pos*1), Y3
	MOVQ    $0, tmp
	VMOVDQU (tbl)(tmp*1), X0
    VMOVDQU 16(tbl)(tmp*1), X1
    VINSERTI128  $1, X0, Y0, Y0
    VINSERTI128  $1, X1, Y1, Y1
    VPSRLQ  $4, Y3, Y4
    VPAND   Y2, Y4, Y4
    VPAND   Y2, Y3, Y5
    VPSHUFB Y4, Y1, Y3
    VPSHUFB Y5, Y0, Y6
    VPXOR   Y3, Y6, Y3
    MOVQ   $0, tmp1
    MOVQ   (parity)(tmp1*1), tmp2
    VPXOR (tmp2)(pos*1), Y3, Y3
    VMOVDQU Y3, (tmp2)(pos*1)
    MOVQ   vect, tmp3
    SUBQ   $2, tmp3
    JL ret
next_vect:
	ADDQ   $32, tmp
	VMOVDQU (tbl)(tmp*1), X0
    VMOVDQU 16(tbl)(tmp*1), X1
    VINSERTI128  $1, X0, Y0, Y0
    VINSERTI128  $1, X1, Y1, Y1
	VPSHUFB  Y4, Y1, Y3
	VPSHUFB  Y5, Y0, Y6
	VPXOR Y3, Y6, Y3
	ADDQ   $24, tmp1
	MOVQ   (parity)(tmp1*1), tmp2
	VPXOR (tmp2)(pos*1), Y3, Y3
	VMOVDQU Y3, (tmp2)(pos*1)
	SUBQ   $1, tmp3
	JGE   next_vect

	ADDQ  $32, pos
	CMPQ  len, pos
	JNZ   loop32b
	RET

ret:
	RET
