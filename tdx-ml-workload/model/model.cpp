#include <iostream>
#include <vector>
#include <cstring>

// LinearRegression Model
#include "lin_reg.h"
#include "model.h"

// openssl
#include "openssl/aes.h"
#include "openssl/evp.h"
#include "openssl/err.h"

using namespace std;

// AI Model (Decrypted)
char *aimodelbuffer = NULL;

#define SGX_AESGCM_IV_SIZE  12
#define SGX_AESGCM_MAC_SIZE 16
#define AESGCM_256_KEY_SIZE 32

typedef uint8_t aes_gcm_256bit_key_t[AESGCM_256_KEY_SIZE];
typedef uint8_t sgx_aes_gcm_128bit_tag_t[SGX_AESGCM_MAC_SIZE];

uint32_t aes_256gcm_decrypt(const aes_gcm_256bit_key_t *p_key, const uint8_t *p_src,
				uint32_t src_len, uint8_t *p_dst, const uint8_t *p_iv, uint32_t iv_len,
				const uint8_t *p_aad, uint32_t aad_len, const sgx_aes_gcm_128bit_tag_t *p_in_mac)
{
     uint8_t l_tag[SGX_AESGCM_MAC_SIZE];

     if ((src_len >= INT_MAX) || (aad_len >= INT_MAX) || (p_key == NULL) || ((src_len > 0) && (p_dst == NULL)) || ((src_len > 0) && (p_src == NULL))
	 || (p_in_mac == NULL) || (iv_len != SGX_AESGCM_IV_SIZE) || ((aad_len > 0) && (p_aad == NULL))
	 || (p_iv == NULL) || ((p_src == NULL) && (p_aad == NULL)))
     {
	  return 1;
     }
     int len = 0;
     uint32_t ret = -1;
     EVP_CIPHER_CTX * pState = NULL;

     // Autenthication Tag returned by Decrypt to be compared with Tag created during seal
     memset(&l_tag, 0, SGX_AESGCM_MAC_SIZE);
     memcpy(l_tag, p_in_mac, SGX_AESGCM_MAC_SIZE);

     do {
	  // Create and initialise the context
	  if (!(pState = EVP_CIPHER_CTX_new())) {
	       ret = 2;
	       break;
	  }

	  // Initialise decrypt, key and IV
	  if (!EVP_DecryptInit_ex(pState, EVP_aes_256_gcm(), NULL, (unsigned char*)p_key, p_iv)) {
	       break;
	  }

	  // Provide AAD data if exist
	  if (NULL != p_aad) {
	       if (!EVP_DecryptUpdate(pState, NULL, &len, p_aad, aad_len)) {
		    break;
	       }
	  }

	  // Decrypt message, obtain the plaintext output
	  if (!EVP_DecryptUpdate(pState, p_dst, &len, p_src, src_len)) {
	       break;
	  }

	  // Update expected tag value
	  //
	  if (!EVP_CIPHER_CTX_ctrl(pState, EVP_CTRL_GCM_SET_TAG, SGX_AESGCM_MAC_SIZE, l_tag)) {
	       break;
	  }

	  // Finalise the decryption. A positive return value indicates success,
	  // anything else is a failure - the plaintext is not trustworthy.
	  if (EVP_DecryptFinal_ex(pState, p_dst + len, &len) <= 0) {
	       ret = 3;
	       break;
	  }
	  ret = 0;
     } while (0);

     // Clean up and return
     if (pState != NULL) {
	  EVP_CIPHER_CTX_free(pState);
     }

     memset(&l_tag, 0, SGX_AESGCM_MAC_SIZE);

     return ret;
}

uint32_t decrypt_aes_wrapped_secret(uint8_t* wrappedSecret,
                              uint32_t wrappedSecretSize,
                              uint8_t* dek,
                              uint8_t** plaintext)
{
    /*
      wrappedSecret format :
      <IV:SGX_AESGCM_IV_SIZE><CipherText:n><MAC:SGX_AESGCM_MAC_SIZE>
    */

    uint32_t ret_code= 0;
    std::cout << "Received wrappedSecret of size : " << wrappedSecretSize;

    // Plaintext Output Buffer.
    int plaintext_len = wrappedSecretSize - (SGX_AESGCM_IV_SIZE + SGX_AESGCM_MAC_SIZE);
    *plaintext = (uint8_t *) malloc (plaintext_len);
    if (!plaintext) {
        printf("Plaintext buffer memory allocation failed.");
        return -1;
    }

    // Cipher text
    uint32_t cipher_text_len = wrappedSecretSize - (SGX_AESGCM_IV_SIZE + SGX_AESGCM_MAC_SIZE);
    std::cout << "Cipher Text Length : " << cipher_text_len;

    // Copy of DEK
    aes_gcm_256bit_key_t *sk_key = (aes_gcm_256bit_key_t *)malloc (sizeof(uint8_t) * AESGCM_256_KEY_SIZE);
    if (!sk_key) {
        printf("DEK buffer memory allocation failed.");
        return -1;
    }

    memcpy (sk_key, dek, AESGCM_256_KEY_SIZE);

    // Extract the MAC from the transmitted cipher text
    sgx_aes_gcm_128bit_tag_t mac;
    memcpy (mac, wrappedSecret + SGX_AESGCM_IV_SIZE+ plaintext_len, SGX_AESGCM_MAC_SIZE);

    // IV initialisation
    unsigned char iv[SGX_AESGCM_IV_SIZE];
    memcpy(iv, wrappedSecret, SGX_AESGCM_IV_SIZE);

    ret_code = aes_256gcm_decrypt(sk_key, // Key
				  wrappedSecret + SGX_AESGCM_IV_SIZE, // Cipher text
				  cipher_text_len, //Cipher len
				  *plaintext, // Plaintext
				  iv, // Initialisation vector
				  SGX_AESGCM_IV_SIZE, // IV Length
				  NULL, // AAD Buffer
				  0, // AAD Length
				  &mac); // MAC

      free(sk_key);
      sk_key = NULL;
    if (0 != ret_code) {
        printf("Secret decryption failed!");
        return ret_code;
    }

    printf("Secret unwrapped successfully...");
    return 0;
}

int model_decrypt(uint8_t* wrapped_model,
                    size_t wrapped_model_size,
                    uint8_t* wrapped_dek,
                    size_t wrapped_dek_size,
                    uint8_t* swk)
{
    uint32_t ret_code = 0;
    uint8_t* data = NULL;

    ret_code = decrypt_aes_wrapped_secret(wrapped_dek,
                              wrapped_dek_size,
                              swk,
                              &data);

    if (ret_code != 0) {
        printf("Decryption of DEK failed. Check error code.");
        if (data != NULL) {
          free(data);
          data = NULL;
        }
        return ret_code;
    }

    std::cout << "Successfully decrypted DEK.";

    swk = data;

    ret_code = decrypt_aes_wrapped_secret(wrapped_model,
                              wrapped_model_size,
                              swk,
                              &data);

      free(swk);
      swk = NULL;
    if (ret_code != 0) {
        printf("Decryption of model failed. Check error code.");
        if (data != NULL) {
          free(data);
          data = NULL;
        }
        return ret_code;
    }

    std::cout << "Successfully decrypted Model.";

    aimodelbuffer = (char *)data;

    return 0;
}

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
