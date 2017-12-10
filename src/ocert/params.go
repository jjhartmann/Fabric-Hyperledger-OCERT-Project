/*
 *
 * Copyright 2017 Kewei Shi, Jeremy Hartmann, Tuhin Tiwari and Dharvi Verma
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package ocert

import (
    "github.com/Nik-U/pbc"
)

/*
 * Randomly generate a paired group and the corresponding
 * generator of each group
 */
func GenerateSharedParams() *SharedParams {
    sharedParams := new(SharedParams)
    params := pbc.GenerateF(640)
    pairing := params.NewPairing()
    g1 := pairing.NewG1().Rand()
    g2 := pairing.NewG2().Rand()
    sharedParams.Params = params.String()
    sharedParams.G1 = g1.Bytes()
    sharedParams.G2 = g2.Bytes()
    return sharedParams
}
