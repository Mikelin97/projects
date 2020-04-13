import numpy as np
import skimage as sk
import skimage.io as skio
import time 
import os 

def recursive_align(image1, image2): 
	if image1.shape[0] < 400 or image1.shape[1] < 400:
		minssd = float('inf')
		x_offset = 0 
		y_offset = 0
		saved_image2 = image2
		cropped_image1 = sk.util.crop(image1, (image1.shape[0]//5, image1.shape[1]//5))
		for x in range(-15, 15): 
			for y in range(-15, 15): 
				shifted_image2 = np.roll(np.roll(image2, x, axis=1), y, axis=0) 
				cropped_image2 = sk.util.crop(shifted_image2, (shifted_image2.shape[0]//5, shifted_image2.shape[1]//5))
				ssd = np.sum(np.square(np.subtract(cropped_image1, cropped_image2)))
				if ssd < minssd: 
					minssd = ssd
					x_offset = x 
					y_offset = y
					saved_image2 = shifted_image2
		return saved_image2, x_offset, y_offset
	else: 
		resized_image1 = sk.transform.pyramid_reduce(image1, downscale=2, multichannel=False)
		resized_image2 = sk.transform.pyramid_reduce(image2, downscale=2, multichannel=False)
		new_image2, x_offset, y_offset = recursive_align(resized_image1, resized_image2)

		minssd = float('inf')
		saved_image2 = image2
		cropped_image1 = sk.util.crop(image1, (image1.shape[0]//5, image1.shape[1]//5))
		for x in range(x_offset -8 , x_offset + 8): 
			for y in range(y_offset -8, y_offset + 8): 
				shifted_image2 = np.roll(np.roll(image2, x, axis=1), y, axis=0) 
				cropped_image2 = sk.util.crop(shifted_image2, (shifted_image2.shape[0]//5, shifted_image2.shape[1]//5))
				ssd = np.sum(np.square(np.subtract(cropped_image1, cropped_image2)))
				if ssd < minssd: 
					minssd = ssd
					x_offset = x 
					y_offset = y
					saved_image2 = shifted_image2
		return saved_image2, x_offset, y_offset

def align(photo1, photo2): 
	minssd = float('inf')
	minphoto = photo1
	cur_x = 0
	cur_y = 0 
	photo2 = sk.util.crop(photo2, (photo2.shape[0]//5, photo2.shape[1]//5))
	for x in range(-15, 15): 
		for y in range(-15, 15): 
			new_g = np.roll(photo1, x, axis=1)
			cur_g = np.roll(new_g, y, axis=0)
			use_g = sk.util.crop(cur_g, (cur_g.shape[0]//5, cur_g.shape[1]//5))

			curphoto = np.square(np.subtract(photo2, use_g))
			ssd = np.sum(curphoto)	
			if ssd < minssd: 
				minssd = ssd
				minphoto = cur_g
				cur_x = x 
				cur_y = y
	return minphoto, cur_x, cur_y


def main(): 
	folder = "./images"

	for file in os.listdir(folder): 
		if file.endswith((".jpg", ".tif")): 
			im = skio.imread(file)
			im = sk.img_as_float(im)
			height = np.floor(im.shape[0] / 3.0).astype(np.int)

			b = im[:height]
			g = im[height: 2*height]
			r = im[2*height: 3*height]

			start = time.time()
			ag, x1, y1 = recursive_align(b, g)
			ar, x2, y2 = recursive_align(b, r)
			end = time.time()
			print(file)
			print(end - start)
			print(x1, y1)
			print(x2, y2)
			im_out = np.dstack([ar, ag, b])


			file = file.rstrip('.tif')
			if not file.endswith(".jpg"): 
				file = file + ".jpg"
			# save the image
			fname = './out_path/' + file
			skio.imsave(fname, im_out)

		# display the image
			# skio.imshow(im_out)
			# skio.show() 

if __name__ == "__main__": 
	main()


