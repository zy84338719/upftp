#!/usr/bin/env python3
"""UPFTP Logo Generator"""

from PIL import Image, ImageDraw
import os

def create_logo(size):
    """创建UPFTP Logo"""
    img = Image.new('RGBA', (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # 颜色
    orange = (217, 119, 6, 255)
    white = (255, 255, 255, 255)
    
    # 绘制圆形背景
    margin = int(size * 0.06)
    draw.ellipse([margin, margin, size-margin, size-margin], fill=orange)
    
    # 绘制U字母（简化版）
    center = size // 2
    u_width = int(size * 0.35)
    u_height = int(size * 0.35)
    stroke = int(size * 0.09)
    
    # U的左右竖线
    u_left = center - u_width // 2
    u_right = center + u_width // 2
    u_top = int(size * 0.32)
    u_bottom = int(size * 0.72)
    
    draw.rectangle([u_left, u_top, u_left + stroke, u_bottom], fill=white)
    draw.rectangle([u_right - stroke, u_top, u_right, u_bottom], fill=white)
    draw.rectangle([u_left, u_bottom - stroke, u_right, u_bottom], fill=white)
    
    # 绘制上传箭头
    arrow_size = int(size * 0.15)
    arrow_x = center
    arrow_y = int(size * 0.22)
    
    draw.polygon([
        (arrow_x, arrow_y - arrow_size),
        (arrow_x - arrow_size, arrow_y),
        (arrow_x + arrow_size, arrow_y)
    ], fill=white)
    
    draw.rectangle([
        arrow_x - stroke//2,
        arrow_y,
        arrow_x + stroke//2,
        arrow_y + arrow_size//2
    ], fill=white)
    
    # 装饰点
    dot_r = int(size * 0.025)
    dot_y = int(size * 0.78)
    for offset in [-0.15, 0, 0.15]:
        x = int(center + size * offset)
        draw.ellipse([x-dot_r, dot_y-dot_r, x+dot_r, dot_y+dot_r], fill=white)
    
    return img

def main():
    output_dir = '/tmp/upftp/assets'
    os.makedirs(output_dir, exist_ok=True)
    
    # 生成PNG文件
    for size in [512, 256, 128, 64, 32]:
        logo = create_logo(size)
        path = f'{output_dir}/logo-{size}x{size}.png'
        logo.save(path, 'PNG')
        print(f'✓ {path}')
    
    # 生成favicon.ico
    favicon_imgs = [create_logo(s) for s in [16, 32, 48, 64]]
    favicon_path = f'{output_dir}/favicon.ico'
    favicon_imgs[0].save(
        favicon_path,
        format='ICO',
        sizes=[(s, s) for s in [16, 32, 48, 64]],
        append_images=favicon_imgs[1:]
    )
    print(f'✓ {favicon_path}')
    
    # 复制到docs目录
    create_logo(512).save('/tmp/upftp/docs/logo.png', 'PNG')
    print('✓ /tmp/upftp/docs/logo.png')
    
    print('\n✅ 所有Logo已生成！')

if __name__ == '__main__':
    main()
