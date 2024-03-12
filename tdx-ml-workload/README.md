# TrustAuthority Demo

## Build

make

## Usage

**Create TDX Key Transfer Policy on KBS**

**Create AES 256 Key on KBS and associate it with Key Transfer Policy**

**Encrypt the datafile using [encryptor](../encryptor)**

**Push the encrypted datafile under /etc/ on TDVM**

**Install [TDX CLI](https://github.com/intel/trustauthority-client-for-go/blob/main/tdx-cli/README.md#installation)**

**Install TrustAuthority Demo App**

Create trustauthority-demo.env file under /tmp/ with below contents :

TRUSTAUTHORITY_API_URL=https://api.trustauthority.intel.com <br>
TRUSTAUTHORITY_API_KEY=<trustauthority api key> <br>
HTTPS_PROXY=<proxy if any> <br>

**./trustauthority-demo-v1.0.0.bin**

## API Spec

### Get attestation token

* **URL**
  `https://<IP>:12780/taa/v1/token`

* **Method:**
  `GET`

* **Headers:**

  `"Accept" : "application/json"`

* **Success Response:**
  * **Code:** 200 <br>
    **Content:**
    ``` json
      {"attestation_token":"eyJhbGciOiJSUzM4NCIsImtpZCI6IjY1OWY4ZDU3OTA1YjY5MTQxYzk5YjY0MjEzNmU3ZjgzMmEyMTdiZTgiLCJ0eXAiOiJKV1QifQ.eyJtcmVuY2xhdmUiOiI5ZDk4OWMwZjk2YWVkYzAzNmU5NmY5YzQ1NTA5Y2NjYTFjMDdlODRlYTg3NzI2NzMzNTAwMDhmMDk0N2UxOTExIiwibXJzaWduZXIiOiJiMzczNzNhYmJkODZlNGY2ZDM3YTczYjMxNmI0Y2Q5OGVkY2U5NjNjMzFiZjljZGU5ZTUzZTljNDVjMzg3MWZlIiwiaXN2cHJvZGlkIjoxLCJzZWFtc3ZuIjowLCJpc3Zzdm4iOjEsImVuY2xhdmVfaGVsZF9kYXRhIjoiQVFBQkFQVUF6U2NFS1BRS2M2WE91anJXUENuTjhFcE5YVVpXc0VHSWErMUc0UFJra3BWOUJIL1NLNDJseFAzbGV5SnVrRER1S290dUN3T1pqelJqOW9oREF4dTExc0FqdjdENkw5dFBFYlJaTmhYTnpWR1BrZ1JzU0c2NlFSQ2NlQW5SWnhLaXN6cE9XNkl0b3JGeWRkN2JoNkJnSm8rYmhQT2h6UkFGZmV0UWFkQUFKcjY0aWdBVVpNS0p5dVVFSi9veGl3K1Vmakh2bE1lZjZrS0J3UElyNWVseW9VOFpXWDIvRUIwNTRiNHArZWJsYTQwall1a1EyZ1VWSHRFbnQ5NFlCeUVWamd1amU3TnZXbVFlTm56QnBmL3FVK2pkMGxsWkowYUc5dnRsbWRvUzlFQzVjc2FURjd0L3RZSzBaNWEyaDJOS1pNS04xeDViajNSbVFWa09BbWZFb2h2ZC9OWnp5KzNjUytSTEgvOXhlNVRMUVJ0eFo1eTQ5Mnl3WFcxekE0ck9qdWxwWFlFbC9xSlIrVlQzU2F0OVlIQm01dXpJTHpaQml4bnF5YzZHbFB1d0lnMkxwc2hUaXNuV2U0WTZ5TGlpWExpeEI0ajZmZjlLVHk5bHlJVFdrdENxNCtNMHJ6QWNqM2c3UzVMWUFidE0vS3d2eEVVOEcxc0hCcjJ1bEE9PSIsInBvbGljeV9pZHMiOltdLCJ0Y2Jfc3RhdHVzIjoiT1VUX09GX0RBVEUiLCJ0ZWUiOiJTR1giLCJ2ZXIiOiIxLjAiLCJleHAiOjE2MzgyMDUzNTYsImlhdCI6MTYzODE5ODEyNiwiaXNzIjoiQVBTIEF0dGVzdGF0aW9uIFRva2VuIElzc3VlciIsIm5iZiI6MTYzODE5ODEyNn0.rpDoDt9kge2mnHt2g3qNDSUZak40M_S050PhmGRW0Xo9unykbqzd3RN5C6wWUZnZKIxbcgbj5ZxIbE3seK2Wz4J0PYXjdSbtqC_xhIm4JDic0N1uUJEQzg_o1EqL_HQVSmKmAUc_q3h2Ec1X-pIUe2NK_IJIsj1sGLQfn0GfgdTAvJunmZrKQ9i4mlGDy7KyLP12q6mkw1CqvrgF6mXVnA0B29dR6EkCb7RmobKAh2UiplrC8WkBbqFrxxFDpo1IHqQqpRgxfr-Lnirhl-e9n2QdxWT9eT_s0rBLyb_wqweqkNx5clft7GPC-DXYhfCfeVxJYlwpxGonwC_qfAE6AvrZv90_VOxAOJUI-roW8Q56XtjIZmPiQrCdqivwSS0d5WlsmOKyWoCxuUQG-Vh2dGm-Vxvur67QYIiMb82R2sdZyLacNT6F5ht37bmcTp_Wz9AcHFVDqMmDewBh5M28Zrj2vv2EphpfS7FMBVaQi9SV1qW4D7RCvQul6evQSFav"}
    ```

### Get decryption key
Client should send the AttestationToken from the previous step and the "Key Transfer URL"

* **URL**
  `https://<IP>:12780/taa/v1/key`

* **Method:**
  `POST`

* **Headers:**

  `"Accept" : "application/json"` <br>
  `"Content-Type" : "application/json"` <br>

* **Data Params:**

  `attestation_token=[string]` <br>
  `key_transfer_url=[string]` <br>

Request JSON :
```json
{
    "attestation_token" : "eyJhbGciOiJSUzM4NCIsImtpZCI6IjY1OWY4ZDU3OTA1YjY5MTQxYzk5YjY0MjEzNmU3ZjgzMmEyMTdiZTgiLCJ0eXAiOiJKV1QifQ.eyJtcmVuY2xhdmUiOiI5ZDk4OWMwZjk2YWVkYzAzNmU5NmY5YzQ1NTA5Y2NjYTFjMDdlODRlYTg3NzI2NzMzNTAwMDhmMDk0N2UxOTExIiwibXJzaWduZXIiOiJiMzczNzNhYmJkODZlNGY2ZDM3YTczYjMxNmI0Y2Q5OGVkY2U5NjNjMzFiZjljZGU5ZTUzZTljNDVjMzg3MWZlIiwiaXN2cHJvZGlkIjoxLCJzZWFtc3ZuIjowLCJpc3Zzdm4iOjEsImVuY2xhdmVfaGVsZF9kYXRhIjoiQVFBQkFQVUF6U2NFS1BRS2M2WE91anJXUENuTjhFcE5YVVpXc0VHSWErMUc0UFJra3BWOUJIL1NLNDJseFAzbGV5SnVrRER1S290dUN3T1pqelJqOW9oREF4dTExc0FqdjdENkw5dFBFYlJaTmhYTnpWR1BrZ1JzU0c2NlFSQ2NlQW5SWnhLaXN6cE9XNkl0b3JGeWRkN2JoNkJnSm8rYmhQT2h6UkFGZmV0UWFkQUFKcjY0aWdBVVpNS0p5dVVFSi9veGl3K1Vmakh2bE1lZjZrS0J3UElyNWVseW9VOFpXWDIvRUIwNTRiNHArZWJsYTQwall1a1EyZ1VWSHRFbnQ5NFlCeUVWamd1amU3TnZXbVFlTm56QnBmL3FVK2pkMGxsWkowYUc5dnRsbWRvUzlFQzVjc2FURjd0L3RZSzBaNWEyaDJOS1pNS04xeDViajNSbVFWa09BbWZFb2h2ZC9OWnp5KzNjUytSTEgvOXhlNVRMUVJ0eFo1eTQ5Mnl3WFcxekE0ck9qdWxwWFlFbC9xSlIrVlQzU2F0OVlIQm01dXpJTHpaQml4bnF5YzZHbFB1d0lnMkxwc2hUaXNuV2U0WTZ5TGlpWExpeEI0ajZmZjlLVHk5bHlJVFdrdENxNCtNMHJ6QWNqM2c3UzVMWUFidE0vS3d2eEVVOEcxc0hCcjJ1bEE9PSIsInBvbGljeV9pZHMiOltdLCJ0Y2Jfc3RhdHVzIjoiT1VUX09GX0RBVEUiLCJ0ZWUiOiJTR1giLCJ2ZXIiOiIxLjAiLCJleHAiOjE2MzgyMDUzNTYsImlhdCI6MTYzODE5ODEyNiwiaXNzIjoiQVBTIEF0dGVzdGF0aW9uIFRva2VuIElzc3VlciIsIm5iZiI6MTYzODE5ODEyNn0.rpDoDt9kge2mnHt2g3qNDSUZak40M_S050PhmGRW0Xo9unykbqzd3RN5C6wWUZnZKIxbcgbj5ZxIbE3seK2Wz4J0PYXjdSbtqC_xhIm4JDic0N1uUJEQzg_o1EqL_HQVSmKmAUc_q3h2Ec1X-pIUe2NK_IJIsj1sGLQfn0GfgdTAvJunmZrKQ9i4mlGDy7KyLP12q6mkw1CqvrgF6mXVnA0B29dR6EkCb7RmobKAh2UiplrC8WkBbqFrxxFDpo1IHqQqpRgxfr-Lnirhl-e9n2QdxWT9eT_s0rBLyb_wqweqkNx5clft7GPC-DXYhfCfeVxJYlwpxGonwC_qfAE6AvrZv90_VOxAOJUI-roW8Q56XtjIZmPiQrCdqivwSS0d5WlsmOKyWoCxuUQG-Vh2dGm-Vxvur67QYIiMb82R2sdZyLacNT6F5ht37bmcTp_Wz9AcHFVDqMmDewBh5M28Zrj2vv2EphpfS7FMBVaQi9SV1qW4D7RCvQul6evQSFav",
    "key_transfer_url" : "https://foobar.com/kbs/v1/keys/<keyid>/transfer"
}
```

* **Success Response:**
  * **Code:** 200 <br>
    **Content:**
    ``` json
    {
        "wrapped_key":"DAAAABAAAAAwAAAAyhw9koWa+qZfzrBW0e9MDk7BEnfrVJby/aPLfEP5tHwpQLQmsDSV27NBpbqMEu33NvqlurhTCEiLedoD",
        "wrapped_swk":"t+HutgD5OcQfbhp0kG47bTjjce+RcjDqB1r38wIAJ/vWkioRsuSel2gOm52pV87DGLzbkQ2BvJGd1+RTE7bUlCgYs9YZt7Sk8tMGN0O3sXK9NOd+Ms9BOrhsUSwbFqWilftHcdOmkJgvHOx6p/kighDJARLV9UbT0fVj04UoaYWXmptf5OJNOrtLBsO5iy5TeQv+jzyIXRgAU98sFNQyaF1g02RcohgPbYa8wmFShXZ0PWM/Qyu4+D6gQRNsDzKqxXAwhOHrg1AKY/0A5V8ZqnrpUAcTO5/LmLiIXoXILUztHofY0Z3VfTsAj/PCBfpEGyNC6duEoV3Gv/iyazElirO3i7QERPWByAa6W7SixGcsETrew2MxYJiEMQ0iBdZBE4xkpO2RLE1I1qYDlATNwBnaPsqG5vUAVLdISVKDo+2YVIRqBmb9FlxjLVLIeExeTvDIfudLaHoocyblkSPayXLek/JPJl7tuqoYf1k5axmSn9MzlwAB8IFLuJ7mW1TH"
    }
    ```

### Decrypt model
Client should send the WrappedKey and WrappedSwk from the previous step

* **URL**
  `https://<IP>:12780/taa/v1/decrypt`

* **Method:**
  `POST`

* **Headers:**

  `"Content-Type" : "application/json"`

* **Data Params:**

  `wrapped_key=[string]` <br>
  `wrapped_swk=[string]` <br>

Request JSON :
```json
{
    "wrapped_key": "DAAAABAAAAAwAAAAyhw9koWa+qZfzrBW0e9MDk7BEnfrVJby/aPLfEP5tHwpQLQmsDSV27NBpbqMEu33NvqlurhTCEiLedoD",
    "wrapped_swk": "t+HutgD5OcQfbhp0kG47bTjjce+RcjDqB1r38wIAJ/vWkioRsuSel2gOm52pV87DGLzbkQ2BvJGd1+RTE7bUlCgYs9YZt7Sk8tMGN0O3sXK9NOd+Ms9BOrhsUSwbFqWilftHcdOmkJgvHOx6p/kighDJARLV9UbT0fVj04UoaYWXmptf5OJNOrtLBsO5iy5TeQv+jzyIXRgAU98sFNQyaF1g02RcohgPbYa8wmFShXZ0PWM/Qyu4+D6gQRNsDzKqxXAwhOHrg1AKY/0A5V8ZqnrpUAcTO5/LmLiIXoXILUztHofY0Z3VfTsAj/PCBfpEGyNC6duEoV3Gv/iyazElirO3i7QERPWByAa6W7SixGcsETrew2MxYJiEMQ0iBdZBE4xkpO2RLE1I1qYDlATNwBnaPsqG5vUAVLdISVKDo+2YVIRqBmb9FlxjLVLIeExeTvDIfudLaHoocyblkSPayXLek/JPJl7tuqoYf1k5axmSn9MzlwAB8IFLuJ7mW1TH"
}
```

* **Success Response:**
  * **Code:** 200 <br>

### Execute model

* **URL**
  `https://<IP>:12780/taa/v1/execute`

* **Method:**
  `POST`

* **Headers:**

  `"Accept" : "application/json"` <br>
  `"Content-Type" : "application/json"` <br>

* **Data Params:**

  `pregnancies=[string]` <br>
  `blood-glucose=[string]` <br>
  `blood-pressure=[string]` <br>
  `skin-thickness=[string]` <br>
  `insulin=[string]` <br>
  `bmi=[string]` <br>
  `dbf=[string]` <br>
  `age=[string]` <br>

Request JSON :
```json
{
    "pregnancies": "3",
    "blood-glucose": "130",
    "blood-pressure": "78",
    "skin-thickness": "23",
    "insulin": "79",
    "bmi": "28.4",
    "dbf": "0.323",
    "age": "34"
}
```

* **Success Response:**
  * **Code:** 200 <br>
    **Content:**
    ``` json
      {"high-risk":0}
    ```