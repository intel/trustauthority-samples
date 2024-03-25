#ifndef _MODEL_H_
#define _MODEL_H_

// ML Model (Decrypted)
extern char *aimodelbuffer;

#ifdef __cplusplus
extern "C" {
#endif

int model_reset();

int model_predict(double p, double g, double bp, double sk,
		   double in, double bmi, double dbf, double age,
		   int *prediction);

#ifdef __cplusplus
}
#endif

#endif
