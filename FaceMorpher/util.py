from cpselect.cpselect import cpselect
import skimage as sk 
import skimage.io as skio
import json
from scipy.spatial import Delaunay
import matplotlib.pyplot as plt
import numpy as np
from skimage.draw import polygon
from skimage.draw import polygon2mask
from skimage.transform import resize
import imageio
import os

def selPoints(img1, img2, file): 
    '''
    the function is used to select the correspondance points 
    between two aligned images. The last variable is the file 
    to save as. 
    e.g. selPoints("original/haohan_1.jpg", "original/rose_1.jpg", 'points_1.json')
    '''
    controlpointlist = cpselect(img1, img2)
    with open(file, 'w') as f:
        json.dump(controlpointlist, f)

def preProcessPoints(file): 
    '''
    After the correspondance points are picked for two images, 
    the points are read from the json file. 
    Return: two numpy arrays containing the points for each image. 
    e.g. img1Co, img2Co = preProcessPoints('points_1.json')
    '''
    img1 = []
    img2 = []
    with open(file, 'r') as f: 
        datastore = json.load(f)
    
    for point in datastore: 
        temp1 = [point['img1_x'], point['img1_y']]
        temp2 = [point['img2_x'], point['img2_y']]
        img1.append(temp1)
        img2.append(temp2)
    img1 = np.array(img1)
    img2 = np.array(img2)
    return img1, img2

def preProcessPointsImm(imName):
    '''
    imName: e.g. 01-1m.jpg
    Return: a list of jpg images in skio.imread format, 
    a numpy array of the correspondance points in the image. 
    ''' 
    absPos = './data/imm_face_db/' + imName
    
    img = skio.imread(absPos)
    width, height, _ = img.shape
    absAsfPos = './data/imm_face_db/' + imName.split('.')[0] + '.asf'
    
    with open(absAsfPos, 'r') as f: 
        data = f.read()
    data.split('\n')
    
    xyIndex = data.split('\n')[16:74]
    points = [[row.split('\t')[2],row.split('\t')[3]] for row in xyIndex]
    result = [[float(point[0]) * height, float(point[1]) * width]for point in points]
    result.extend([[0,0],[width-1,0],[0,height-1],[width-1,height-1]])
    return img, np.array(result)


def triStack(triangle): 
    '''
    given three pairs of (x, y) points (a triangle), 
    triStack turn the vector into a (6, 6) matrix. 
    of the form: 
    [x1, y1, 1, 0, 0, 0]
    [0, 0, 0, x1, y1, 1]
    [x2, y2, 1, 0, 0, 0]
    [0, 0, 0, x2, y2, 1]
    [x3, y3, 1, 0, 0, 0]
    [0, 0, 0, x3, y3, 1]
    '''
    triangle = triangle[:, (1, 0)]
    row1 = np.hstack((triangle[0], [1]))
    row1New = np.hstack((row1 ,[0, 0, 0]))
    row2 = np.hstack(([0, 0, 0], row1))
    row3 = np.hstack((triangle[1], [1]))
    row3New = np.hstack((row3 ,[0, 0, 0]))
    row4 = np.hstack(([0, 0, 0], row3))
    row5 = np.hstack((triangle[2], [1]))
    row5New = np.hstack((row5 ,[0, 0, 0]))
    row6 = np.hstack(([0, 0, 0], row5))
    result = np.vstack((row1New, row2, row3New, row4, row5New, row6))
    return result 


def triStack1(triangle): 
    '''
    flatten the three pairs of (x, y)
    [x1, y1, x2, y2, x3, y3]
    '''
    triangle = triangle[:, (1, 0)]
    a = np.hstack((triangle[0], triangle[1], triangle[2]))
    return a

def computeAffine(tri1_pts, tri2_pts):
    '''
    given two triangles with different shapes, 
    the affine transformation matrix is calculated. 
    Return: the inverse of the affine matrix
    '''
    matrix1 = triStack(tri1_pts)
    matrix2 = triStack1(tri2_pts)
    result = np.linalg.solve(matrix1, matrix2)
    result = result.reshape(2, 3)
    result = np.vstack((result, [0, 0, 1]))
    result = np.linalg.inv(result)
    return result 

def MidwayTris(tri1, tri2, warp_frac=0.5): 
    '''
    given two triangles, the intermediate shape 
    is calculated based on the warp_frac. 
    warp_frac is default as 0.5. 
    '''
    return (1 - warp_frac) * tri1 + warp_frac * tri2

def crossDissolve(img1, img2, dis_frac=0.5): 
    '''
    mix the color of two images. 
    '''
    return (1 - dis_frac) * img1 + dis_frac * img2


