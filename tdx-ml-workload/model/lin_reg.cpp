#include <vector>
#include <assert.h>
#include <iostream>
#include <string>

using namespace std;

double dotProduct(std::vector<double> const &a, std::vector<double> const &b)
{
    assert(a.size() == b.size());
    double dot_product = 0;
    for(int i = 0; i < a.size(); i++)
        dot_product += (a.at(i) * b.at(i));

    return dot_product;
}

std::vector<double> normalizeInput (std::vector<double> inputData) {
  inputData[0] = inputData[0] / 28;
  inputData[1] = inputData[1] / 200;
  inputData[2] = inputData[2] / 125;
  inputData[3] = inputData[3] / 100;
  inputData[4] = inputData[4] / 850;
  inputData[5] = inputData[5] / 68;
  inputData[6] = inputData[6] / 2.45;
  inputData[7] = inputData[7] / 100;
  return inputData;
}

std::vector<double> linreg_load_model_frombuffer (std::string string,
						       double *t)
{
  std::vector<std::string> parse_result;
  std::vector<double> result = {0,0,0,0,0,0,0,0};
  std::string temp;
  int markbegin = 0;
  int markend = 0;

  for (int i = 0; i < string.length(); ++i) {
    if (string[i] == ' ') {
      markend = i;
      parse_result.push_back(string.substr(markbegin, markend - markbegin));
      markbegin = (i + 1);
    }
  }

  int i = 0;
  for (i = 0; i < parse_result.size()-1; i++){
    result[i] = stof(parse_result[i]);
  }

  *t = stof(parse_result[i]);

  return result;
}

int linreg_classify (std::vector<double> &x, std::vector<double> &w, double t) {
  int prediction = -1;

  x = normalizeInput (x);

  double dot_product = dotProduct(x,w);
  if(dot_product > t)
    prediction = 1;
  else
    prediction = 0;

  return prediction;
}
