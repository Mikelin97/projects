import numpy as np
from skimage.draw import polygon
from utils import *

def morph(img1, img2, img1Co, img2Co, tri, warp_frac=0.5, dissolve_frac=0.5):  
    haohanTri = img1Co[tri.simplices.copy()]
    roseTri = img2Co[tri.simplices.copy()]
    resultTri = MidwayTris(haohanTri, roseTri, warp_frac)
    morphImg1 = np.zeros(img1.shape)
    morphImg2 = np.zeros(img2.shape)

    height, width, _ = img2.shape

    for x in range(haohanTri.shape[0]):
        
        A = computeAffine(haohanTri[x], resultTri[x])
        B = computeAffine(roseTri[x], resultTri[x])
        x_index = resultTri[x][:, 0]
        y_index = resultTri[x][:, 1]    
        rr, cc = polygon(y_index, x_index)

        test = np.vstack((rr, cc, [1 for x in range(len(rr))]))

        result1 = np.array(np.dot(A, test), dtype=int)
        result2 = np.array(np.dot(B, test), dtype=int)
        morphImg1[np.clip(test[0,:], 0, height -1), np.clip(test[1,:], 0, width-1), :] = \
            img1[np.clip(result1[0,:], 0, height -1), np.clip(result1[1,:], 0, width-1), :]
        morphImg2[np.clip(test[0,:], 0, height -1), np.clip(test[1,:], 0, width-1), :] = \
            img2[np.clip(result2[0,:], 0, height -1), np.clip(result2[1,:], 0, width-1), :]
    
    morphImg1 = morphImg1/255.0
    morphImg2 = morphImg2/255.0
    morphImg = crossDissolve(morphImg1, morphImg2, dissolve_frac)
     
    return morphImg


def morph1Img(img1, img1Co, meanCo, tri):  
    haohanTri = img1Co[tri.simplices.copy()]
    resultTri = meanCo[tri.simplices.copy()]
    morphImg1 = np.zeros(img1.shape)
    width, height, _ = img1.shape

    for x in range(haohanTri.shape[0]):    
        A = computeAffine(haohanTri[x], resultTri[x])
        x_index = resultTri[x][:, 0]
        y_index = resultTri[x][:, 1]    
        rr, cc = polygon(y_index, x_index)

        test = np.vstack((rr, cc, [1 for x in range(len(rr))]))
        test = np.array(test, dtype=int)
        result1 = np.array(np.dot(A, test), dtype=int)
        morphImg1[np.clip(test[0,:], 0, width -1), np.clip(test[1,:], 0, height-1), :] = \
            img1[np.clip(result1[0,:], 0, width -1), np.clip(result1[1,:], 0, height-1), :]
    
     
    return morphImg1