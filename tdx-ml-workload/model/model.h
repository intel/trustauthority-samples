#ifndef _MODEL_H_
#define _MODEL_H_

#include <stdint.h>
#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

int model_reset();

int model_predict(double p, double g, double bp, double sk,
		   double in, double bmi, double dbf, double age,
		   int *prediction);

int model_decrypt(uint8_t* wrapped_model,
    size_t wrapped_model_size,
    uint8_t* wrapped_dek,
    size_t wrapped_dek_size,
    uint8_t* swk);

#ifdef __cplusplus
}
#endif

#endif
