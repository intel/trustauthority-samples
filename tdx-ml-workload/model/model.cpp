#include <iostream>
#include <vector>
#include <cstring>

// LinearRegression Model
#include "lin_reg.h"
#include "model.h"

// AI Model (Decrypted)
char *aimodelbuffer = NULL;

// Use this function to reset any state the model might hold.
int model_reset() {
    if (aimodelbuffer != NULL) {
      free(aimodelbuffer);
      aimodelbuffer = NULL;
      std::cout << "Dropped the decrypted model reference.";
    } else {
      std::cout << "Model is already clean. Nothing to reset.";
    }

    return 0;
}

int model_predict(double pregnancies,
				 double glucose,
				 double bloodpressure,
				 double skinthickness,
				 double insulin,
				 double bmi,
				 double dbf,
				 double age,
				 int *prediction)
{
  // Decrypted AI model buffer.
  if (aimodelbuffer == NULL) {
    return -1;
  }

  // Load x and w from buffer
  std::vector<double> weights;
  double t;

  std::string model(aimodelbuffer);
  weights = linreg_load_model_frombuffer (model, &t);

  std::vector<double> input = {0,0,0,0,0,0,0,0};
  input[0] = pregnancies;
  input[1] = glucose;
  input[2] = bloodpressure;
  input[3] = skinthickness;
  input[4] = insulin;
  input[5] = bmi;
  input[6] = dbf;
  input[7] = age;

  *prediction = linreg_classify (input, weights, t);
  return 0;
}
