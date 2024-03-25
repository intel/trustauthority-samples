#ifndef _LIN_REG_H__
#define _LIN_REG_H__

#include <vector>
#include <string>

int linreg_classify (std::vector<double> &x, std::vector<double> &w, double t);
std::vector<double> linreg_load_model_frombuffer (std::string string,
						  double *t);

#endif
